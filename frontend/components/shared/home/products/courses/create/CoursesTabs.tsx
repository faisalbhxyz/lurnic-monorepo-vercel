"use client";

import React, { Suspense, useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import Basics from "./details/Basics";
import Curriculum from "./curriculum/Curriculum";
import QuizTable from "./quiz-evaluation/QuizTable";
import AssignmentTable from "./assignment-evaluation/AssignmentTable";
import ReviewsTable from "./reviews/ReviewsTable";
import NoticesTable from "./notices/NoticesTable";
import SettingTabs from "./settings/SettingTabs";
import { FormProvider, useForm } from "react-hook-form";
import { CourseSchema, TCourseSchema } from "@/schema/course.schema";
import { zodResolver } from "@hookform/resolvers/zod";
import CreateNewChapterModal from "./curriculum/NewChapterModal";
import Image from "next/image";
import { LuListVideo } from "react-icons/lu";
import Button from "@/components/ui/Button";
import EditChapter from "./curriculum/EditChapter";
import CreateLessonModal from "./curriculum/CreateLessonModal";
import UpdateLessonModal from "./curriculum/UpdateLessonModal";
import { useSession } from "next-auth/react";
import axiosInstance from "@/lib/axiosInstance";
import { toast } from "sonner";
import CreateQuizModal from "./curriculum/CreateQuizModal";
import AddNewQuestion from "./curriculum/AddNewQuestion";
import CreateAssignmentModal from "./curriculum/CreateAssignmentModal";
import UpdateAssignmentModal from "./curriculum/UpdateAssignmentModal";
import UpdateQuizModal from "./curriculum/UpdateQuizModal";
import UpdateQuestion from "./curriculum/UpdateQuestion";
import {
  dbTimeToPickerFormat,
  getFirstFormError,
} from "@/lib/helpers";
import axios from "axios";
import { normalizeChapterItemPositions, buildChapterItems, splitChapterItems } from "@/lib/chapterItems";

const toDbId = (id?: number | null) =>
  id != null && Number(id) > 0 ? Number(id) : undefined;

const DEFAULT_CERTIFICATE_TEMPLATE = "/images/Certificat-14.jpg";

const normalizeCertificateSettingsForForm = (
  cert?: CourseCertificateSettings | null
) => {
  const hasRow = Boolean(cert?.id);
  const countLessons = hasRow ? Boolean(cert?.count_lessons) : true;
  const countQuizzes = hasRow ? Boolean(cert?.count_quizzes) : true;
  const countAssignments = hasRow ? Boolean(cert?.count_assignments) : true;
  const allCountsOff = !countLessons && !countQuizzes && !countAssignments;

  return {
    is_enabled: cert?.is_enabled ?? false,
    completion_percent: cert?.completion_percent ?? 100,
    count_lessons: allCountsOff ? true : countLessons,
    count_quizzes: allCountsOff ? true : countQuizzes,
    count_assignments: allCountsOff ? true : countAssignments,
    template_path: cert?.template_path?.trim() || DEFAULT_CERTIFICATE_TEMPLATE,
    title: cert?.title ?? "Certificate of Completion",
    subtitle_one: cert?.subtitle_one ?? "",
    subtitle_two: cert?.subtitle_two ?? "",
    owner_signature: cert?.owner_signature
      ? { isDBImg: true, name: cert.owner_signature }
      : null,
    instructor_signature: cert?.instructor_signature
      ? { isDBImg: true, name: cert.instructor_signature }
      : null,
  };
};

const serializeLessonSourceForApi = (source?: {
  data?: unknown;
  isFile?: boolean;
  playback_time?: string | null;
}) => ({
  data:
    typeof source?.data === "string"
      ? source.data
      : source?.data != null
        ? String(source.data)
        : "",
  is_file: Boolean(source?.isFile),
  playback_times: source?.playback_time ?? null,
});

const tabForFieldPath = (path: string): string => {
  if (path.startsWith("course_chapters")) return "Curriculum";
  if (
    path.startsWith("general_settings") ||
    path.startsWith("course_instructors") ||
    path.startsWith("certificate_settings")
  ) {
    return "Settings";
  }
  return "Details";
};

const serializeChaptersForApi = (chapters: TCourseSchema["course_chapters"]) =>
  chapters.map((chapter) => {
    const { course_lessons, quizzes, assignments } = splitChapterItems(
      buildChapterItems(chapter)
    );

    return {
    id: toDbId(chapter.id),
    position: chapter.position,
    title: chapter.title,
    description: chapter.description ?? null,
    access: chapter.access,
    course_lessons: course_lessons.map((lesson) => ({
      id: toDbId(lesson.id),
      position: lesson.position ?? 0,
      title: lesson.title,
      description: lesson.description ?? null,
      lesson_type: lesson.lesson_type,
      source_type: lesson.source_type,
      source: serializeLessonSourceForApi(lesson.source),
      is_published: lesson.is_published,
      is_public: lesson.is_public,
      is_scheduled: lesson.is_scheduled,
      schedule_date: lesson.schedule_date ?? null,
      schedule_time: lesson.schedule_time ?? null,
      show_comming_soon: lesson.show_comming_soon,
      resources: lesson.resources ?? null,
    })),
    quizzes: quizzes.map((quiz) => ({
      id: toDbId(quiz.id),
      position: quiz.position ?? 0,
      title: quiz.title,
      instructions: quiz.instructions,
      is_published: quiz.is_published,
      randomize_questions: quiz.randomize_questions,
      single_quiz_view: quiz.single_quiz_view,
      time_limit: quiz.time_limit,
      time_limit_option: quiz.time_limit_option,
      total_visible_questions: quiz.total_visible_questions,
      reveal_answers: quiz.reveal_answers,
      enable_retry: quiz.enable_retry,
      retry_attempts: quiz.retry_attempts,
      minimum_pass_percentage: quiz.minimum_pass_percentage,
      questions: (quiz.questions ?? []).map((question) => ({
        id: toDbId(question.id),
        title: question.title,
        details: question.details ?? null,
        media: question.media ?? null,
        options: question.options ?? null,
        correct_answer: question.correct_answer ?? null,
        type: question.type,
        marks: question.marks,
        answer_required: question.answer_required,
        answer_explanation: question.answer_explanation ?? null,
      })),
    })),
    assignments: assignments.map((assignment) => ({
      id: toDbId(assignment.id),
      position: assignment.position ?? 0,
      title: assignment.title,
      instructions: assignment.instructions,
      attachments:
        (assignment.attachments ?? []).filter(
          (item) => !(item instanceof File)
        ) ?? null,
      is_published: assignment.is_published,
      time_limit: assignment.time_limit,
      time_limit_option: assignment.time_limit_option,
      file_upload_limit: assignment.file_upload_limit,
      total_marks: assignment.total_marks,
      minimum_pass_marks: assignment.minimum_pass_marks,
    })),
  };
  });

const tabs = [
  { id: 1, label: "Details" },
  { id: 2, label: "Curriculum" },
  { id: 3, label: "Quizzes" },
  { id: 4, label: "Assignments" },
  // { id: 5, label: "Reviews" },
  // { id: 6, label: "Notice" },
  { id: 7, label: "Settings" },
];

export default function CoursesTabs({
  categories,
  subcategories,
  instructors,
  isEdit = false,
  courseDetails,
}: {
  categories: ICategory[] | null;
  subcategories: ISubCategory[] | null;
  instructors: IInstructor[] | null;
  isEdit?: boolean;
  courseDetails?: CourseDetails | null;
}) {
  const router = useRouter();
  const { data: session } = useSession();
  const searchParams = useSearchParams();

  const defaultTab = tabs[0].label;
  const tabFromQuery = searchParams.get("tab");
  const isValidTab = tabs.some((tab) => tab.label === tabFromQuery);

  const [activeTab, setActiveTab] = useState(
    isValidTab ? tabFromQuery! : defaultTab
  );

  // Update URL when tab changes
  const handleTabChange = (tabLabel: string) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set("tab", tabLabel);
    router.replace(`?${params.toString()}`);
    setActiveTab(tabLabel);
  };

  // Keep URL changes in sync (in case user uses browser navigation)
  useEffect(() => {
    if (tabFromQuery && tabFromQuery !== activeTab && isValidTab) {
      setActiveTab(tabFromQuery);
    }
  }, [tabFromQuery]);

  const formMethods = useForm<TCourseSchema>({
    resolver: zodResolver(CourseSchema),
    defaultValues: {
      author_id: 1,
      visibility: "public",
      is_scheduled: false,
      schedule_date: null,
      schedule_time: null,
      pricing_model: "free",
      regular_price: 0,
      sale_price: 0,
      show_comming_soon: false,
      course_chapters: [
        {
          _id: Date.now(),
          title: "New Chapter",
          access: "draft",
          position: 0,
          description: "",
          course_lessons: [],
          quizzes: [],
          assignments: [],
        },
      ],
      general_settings: {
        difficulty_level: "all",
        maximum_student: 0,
        language: "english",
        duration: null,
        sub_category_id: null,
      },
      certificate_settings: {
        is_enabled: false,
        completion_percent: 100,
        count_lessons: true,
        count_quizzes: true,
        count_assignments: true,
        template_path: DEFAULT_CERTIFICATE_TEMPLATE,
        title: "Certificate of Completion",
        subtitle_one: "",
        subtitle_two: "",
        owner_signature: null,
        instructor_signature: null,
      },
    },
  });

  useEffect(() => {
    if (isEdit && courseDetails) {
      console.log("courseDetails", courseDetails);

      const instructors = courseDetails.course_instructors.map(
        (ins) => ins.instructor.id
      );

      const chapters = courseDetails.course_chapters
        .slice()
        .sort((a, b) => (a.position ?? 0) - (b.position ?? 0))
        .map((chapter) => {
        return {
          ...chapter,
          id: chapter.id,
          _id: chapter.id,
          course_lessons: (chapter.course_lessons ?? [])
            .map((lesson) => {
            const source = lesson.source || {};
            return {
              ...lesson,
              id: lesson.id,
              _id: lesson.id,
              type: "lesson" as const,
              source: {
                data: source.data?.data ?? source.data,
                playback_time: source.data?.playback_times,
                isFile: false,
              },
            };
          })
            .sort(
              (a, b) =>
                ((a as { position?: number }).position ?? 0) -
                ((b as { position?: number }).position ?? 0)
            ),
          assignments: (chapter.assignments ?? []).map((assignment) => {
            return {
              ...assignment,
              id: assignment.id,
              _id: assignment.id,
              type: "assignment",
              position: (assignment as { position?: number }).position,
            };
          }),
          quizzes: (chapter.quizzes ?? []).map((quiz) => {
            return {
              ...quiz,
              id: quiz.id,
              _id: quiz.id,
              type: "quiz",
              position: (quiz as { position?: number }).position,
              questions: (quiz.questions ?? []).map((q) => {
                return {
                  ...q,
                  id: q.id,
                  _id: q.id,
                };
              }),
            };
          }),
        };
      });

      formMethods.reset({
        title: courseDetails.title,
        summary: courseDetails.summary,
        description: courseDetails.description,
        overview: courseDetails.overview,
        visibility: courseDetails.visibility,
        intro_video: courseDetails.intro_video && {
          type: (courseDetails.intro_video?.type as any) || "youtube",
          source: courseDetails.intro_video.source
            ? courseDetails.intro_video.source
            : null,
        },
        is_scheduled: courseDetails.is_scheduled ? true : false,
        schedule_date: courseDetails.schedule_date
          ? (courseDetails.schedule_date as any)
          : null,
        schedule_time: courseDetails.schedule_time
          ? dbTimeToPickerFormat(courseDetails.schedule_time) // pre-select time
          : "",
        featured_image: {
          isDBImg: true,
          name: courseDetails.featured_image,
          size: 13495,
          type: "image/jpeg",
        },
        pricing_model: courseDetails.pricing_model,
        regular_price: courseDetails.regular_price || 0,
        sale_price: courseDetails.sale_price || 0,
        show_comming_soon: courseDetails?.show_comming_soon || false,
        tags: courseDetails.tags,
        course_chapters: chapters.map((chapter) =>
          normalizeChapterItemPositions({
          ...chapter,
          course_lessons: chapter.course_lessons?.map((lesson) => ({
            ...lesson,
            is_scheduled: lesson.is_scheduled ? true : false, // update as needed
            schedule_date: lesson.schedule_date ?? null, // update as needed
            schedule_time: lesson.schedule_time
              ? dbTimeToPickerFormat(lesson.schedule_time)
              : "",
            show_comming_soon: lesson.show_comming_soon || false,
            resources:
              // @ts-ignore
              lesson.resources?.map((r) => ({
                id: r.id,
                course_id: r.course_id,
                isDBImg: true,
                name: r.name || r.title,
                url: r.url || r.file_path,
                size: r.size || 13495,
                type: r.type || r.mime_type || "application/octet-stream",
              })) || null,
          })),
        })
        ) as any,
        course_instructors: instructors as any,
        author_id: courseDetails.author_id,
        general_settings: {
          difficulty_level: courseDetails.general_settings
            .difficulty_level as any,
          maximum_student: courseDetails.general_settings
            .maximum_student as any,
          language: courseDetails.general_settings.language,
          duration: courseDetails.general_settings.duration,
          category_id: courseDetails.general_settings.category_id,
          sub_category_id:
            courseDetails.general_settings.sub_category_id &&
            courseDetails.general_settings.sub_category_id > 0
              ? courseDetails.general_settings.sub_category_id
              : null,
        },
        certificate_settings: normalizeCertificateSettingsForForm(
          courseDetails.certificate_settings
        ),
      });
    }
  }, [isEdit, courseDetails]);

  const renderContent = () => {
    switch (activeTab) {
      case "Details":
        return <Basics />;
      case "Curriculum":
        return <Curriculum />;
      case "Quizzes":
        return courseDetails?.id ? (
          <div className="border p-5 rounded-lg mt-5">
            <QuizTable courseId={courseDetails.id} />
          </div>
        ) : (
          <div className="border p-5 rounded-lg mt-5 text-sm text-gray-500">
            Save the course first to view quiz submissions.
          </div>
        );
      case "Assignments":
        return courseDetails?.id ? (
          <div className="border p-5 rounded-lg mt-5">
            <AssignmentTable courseId={courseDetails.id} />
          </div>
        ) : (
          <div className="border p-5 rounded-lg mt-5 text-sm text-gray-500">
            Save the course first to view assignment submissions.
          </div>
        );
      // case "Reviews":
      //   return (
      //     <div className="border p-5 rounded-lg mt-5">
      //       <ReviewsTable />
      //     </div>
      //   );
      // case "Notice":
      //   return (
      //     <div className="border p-5 rounded-lg mt-5">
      //       <NoticesTable />
      //     </div>
      //   );
      case "Settings":
        return (
          <div className="border p-5 rounded-lg mt-5">
            <Suspense fallback={<p>Loading...</p>}>
              <SettingTabs
                categories={categories}
                subcategories={subcategories}
                instructors={instructors}
              />
            </Suspense>
          </div>
        );
      default:
        return null;
    }
  };

  const [loading, setLoading] = useState(false);

  const getErrorMessage = (err: unknown): string => {
    if (axios.isAxiosError(err)) {
      const data = err.response?.data as
        | { error?: unknown; message?: unknown }
        | undefined;
      const msg =
        (typeof data?.error === "string" && data.error.trim()) ||
        (typeof data?.message === "string" && data.message.trim());
      if (msg) return msg;

      const status = err.response?.status;
      if (status === 401) return "Unauthorized. Please login again.";
      if (status === 413) return "Upload too large.";
      if (status) return `Request failed (HTTP ${status}).`;

      const code = err.code;
      if (code === "ECONNABORTED") return "Request timed out. Try again.";
      if (code === "ERR_NETWORK") return "Network error. Check connection/API.";
    }

    if (err instanceof Error && err.message.trim()) return err.message;
    return "Something went wrong.";
  };

  const handleInvalid = (errors: typeof formMethods.formState.errors) => {
    const first = getFirstFormError(errors as Record<string, unknown>);
    if (first) {
      toast.error(`${first.message} (${first.path})`);
      handleTabChange(tabForFieldPath(first.path));
      return;
    }
    toast.error("Please fix the highlighted form errors before saving.");
  };

  const handleSave = (data: TCourseSchema) => {
    if (!session?.accessToken) {
      toast.error("Session expired. Please login again.");
      router.push("/login");
      return;
    }

    console.log(data);
    // #region agent log
    fetch("/api/debug-log", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Debug-Session-Id": "01d620",
      },
      body: JSON.stringify({
        sessionId: "01d620",
        runId: "pre-fix",
        hypothesisId: "H1",
        location: "CoursesTabs.tsx:handleSave(entry)",
        message: "Course save submit payload (schedule fields)",
        data: {
          isEdit,
          hasAccessToken: Boolean(session?.accessToken),
          is_scheduled: data.is_scheduled,
          schedule_date: data.schedule_date,
          schedule_time: data.schedule_time,
          visibility: data.visibility,
          show_comming_soon: data.show_comming_soon,
        },
        timestamp: Date.now(),
      }),
    }).catch(() => {});
    // #endregion agent log
    if (loading) return;
    setLoading(true);
    let fd: FormData;
    try {
      fd = new FormData();
      fd.append("title", data.title);
      fd.append("summary", data.summary);
      if (data.description) {
        fd.append("description", data.description);
      }
      fd.append("visibility", data.visibility);
      fd.append("is_scheduled", String(data.is_scheduled));
      if (data.is_scheduled) {
        fd.append("schedule_date", String(data.schedule_date));
        fd.append("schedule_time", String(data.schedule_time));
        fd.append("show_comming_soon", String(data.show_comming_soon));
      }
      fd.append("pricing_model", data.pricing_model);
      fd.append("regular_price", String(data.regular_price));
      fd.append("sale_price", String(data.sale_price));
      if (data.featured_image && !data.featured_image.isDBImg) {
        fd.append("featured_image", data.featured_image);
      }
      if (data.intro_video) {
        fd.append("intro_video", JSON.stringify(data.intro_video));
      }
      fd.append("tags", JSON.stringify(data.tags || []));
      fd.append("author_id", String(data.author_id));
      fd.append("overview", JSON.stringify(data.overview));
      fd.append("course_chapters", JSON.stringify(serializeChaptersForApi(data.course_chapters)));
      fd.append("general_settings", JSON.stringify(data.general_settings));
      fd.append(
        "certificate_settings",
        JSON.stringify({
          is_enabled: data.certificate_settings.is_enabled,
          completion_percent: data.certificate_settings.completion_percent,
          count_lessons: data.certificate_settings.count_lessons,
          count_quizzes: data.certificate_settings.count_quizzes,
          count_assignments: data.certificate_settings.count_assignments,
          template_path: data.certificate_settings.template_path,
          title: data.certificate_settings.title || null,
          subtitle_one: data.certificate_settings.subtitle_one || null,
          subtitle_two: data.certificate_settings.subtitle_two || null,
        })
      );
      if (
        data.certificate_settings.owner_signature &&
        !data.certificate_settings.owner_signature.isDBImg &&
        data.certificate_settings.owner_signature instanceof File
      ) {
        fd.append("owner_signature", data.certificate_settings.owner_signature);
      }
      if (
        data.certificate_settings.instructor_signature &&
        !data.certificate_settings.instructor_signature.isDBImg &&
        data.certificate_settings.instructor_signature instanceof File
      ) {
        fd.append(
          "instructor_signature",
          data.certificate_settings.instructor_signature
        );
      }
      fd.append("course_instructors", JSON.stringify(data.course_instructors));

      // Only append actual new uploads; ignore DB-backed resource objects.
      data.course_chapters.forEach((chapter, chapterIndex) => {
        chapter.course_lessons?.forEach((lesson, lessonIndex) => {
          lesson.resources?.forEach((file) => {
            if (file instanceof File) {
              fd.append(
                `resources[${chapterIndex}][${lessonIndex}][]`,
                file,
                file.name
              );
            }
          });
        });
        chapter.assignments?.forEach((assignment, assignmentIndex) => {
          assignment.attachments?.forEach((file) => {
            if (file instanceof File) {
              fd.append(
                `assignment_attachments[${chapterIndex}][${assignmentIndex}][]`,
                file,
                file.name
              );
            }
          });
        });
      });
    } catch (e: unknown) {
      console.log("[ERROR] build FormData", e);
      toast.error(getErrorMessage(e));
      setLoading(false);
      return;
    }

    if (isEdit && courseDetails) {
      // console.log("[EDIT]", data);
      axiosInstance
        .put(`/private/course/update/${courseDetails.id}`, fd, {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: `Bearer ${session.accessToken}`,
          },
        })
        .then((res) => {
          // #region agent log
          fetch("/api/debug-log", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                "X-Debug-Session-Id": "01d620",
              },
              body: JSON.stringify({
                sessionId: "01d620",
                runId: "pre-fix",
                hypothesisId: "H1",
                location: "CoursesTabs.tsx:handleSave(edit-success)",
                message: "Course edit success",
                data: {
                  status: res.status,
                  ok: Boolean(res.data?.success ?? res.data?.ok),
                  message: res.data?.message,
                },
                timestamp: Date.now(),
              }),
            }).catch(() => {});
          // #endregion agent log
          toast.success(res.data.message);
          router.push(`/courses`);
        })
        .catch((error) => {
          console.log("[ERROR]", error);
          // #region agent log
          fetch("/api/debug-log", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                "X-Debug-Session-Id": "01d620",
              },
              body: JSON.stringify({
                sessionId: "01d620",
                runId: "pre-fix",
                hypothesisId: "H1",
                location: "CoursesTabs.tsx:handleSave(edit-error)",
                message: "Course edit failed",
                data: {
                  status: error?.response?.status,
                  error: error?.response?.data?.error ?? "unknown",
                },
                timestamp: Date.now(),
              }),
            }).catch(() => {});
          // #endregion agent log
          toast.error(getErrorMessage(error));
        })
        .finally(() => setLoading(false));
    } else {
      axiosInstance
        .post("/private/course/create", fd, {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: `Bearer ${session.accessToken}`,
          },
        })
        .then((res) => {
          // #region agent log
          fetch("/api/debug-log", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                "X-Debug-Session-Id": "01d620",
              },
              body: JSON.stringify({
                sessionId: "01d620",
                runId: "pre-fix",
                hypothesisId: "H1",
                location: "CoursesTabs.tsx:handleSave(create-success)",
                message: "Course create success",
                data: {
                  status: res.status,
                  ok: Boolean(res.data?.success ?? res.data?.ok),
                  message: res.data?.message,
                },
                timestamp: Date.now(),
              }),
            }).catch(() => {});
          // #endregion agent log
          toast.success(res.data.message);
          router.push(`/courses`);
        })
        .catch((error) => {
          console.log("[ERROR]", error);
          // #region agent log
          fetch("/api/debug-log", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                "X-Debug-Session-Id": "01d620",
              },
              body: JSON.stringify({
                sessionId: "01d620",
                runId: "pre-fix",
                hypothesisId: "H1",
                location: "CoursesTabs.tsx:handleSave(create-error)",
                message: "Course create failed",
                data: {
                  status: error?.response?.status,
                  error: error?.response?.data?.error ?? "unknown",
                },
                timestamp: Date.now(),
              }),
            }).catch(() => {});
          // #endregion agent log
          toast.error(getErrorMessage(error));
        })
        .finally(() => setLoading(false));
    }
  };

  return (
    <>
      <div className="border-b border-gray-300 flex space-x-4 mb-4">
        {tabs.map((tab) => (
          <button
            type="button"
            key={tab.id}
            onClick={() => handleTabChange(tab.label)}
            className={`px-3 py-2 border-b-2 font-medium text-sm transition-all focus:outline-none ${
              activeTab === tab.label
                ? "border-primary text-primary"
                : "border-transparent text-gray-500 hover:text-gray-700"
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>
      <FormProvider {...formMethods}>
        {/* {JSON.stringify(formMethods.formState.errors, null, 2)} */}
        <div>
          <div className="flex items-center justify-between my-5">
            <div className="flex items-center gap-5">
              <Image
                src={"/images/no-image.png"}
                alt={"image"}
                width={130}
                height={130}
              />
              <div>
                <p className="font-semibold mb-3">Course</p>

                <div className="flex items-center gap-2 text-gray-500 text-sm">
                  <LuListVideo />
                  <p>Recorded Course</p>
                </div>
              </div>
            </div>
            <div className="flex items-center gap-5">
              {/* <button
                type="button"
                className="border px-4 text-sm py-2 rounded-full text-gray-600 font-medium"
              >
                Preview
              </button> */}
              <Button
                type="button"
                disabled={loading}
                onClick={formMethods.handleSubmit(handleSave, handleInvalid)}
                className="disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isEdit ? "Update Course" : "Publish Course"}
              </Button>
            </div>
          </div>
          {renderContent()}
          <CreateNewChapterModal />
          <EditChapter />
          <CreateLessonModal />
          <UpdateLessonModal />
          <CreateQuizModal />
          <AddNewQuestion />
          <CreateAssignmentModal />
          <UpdateAssignmentModal />
          <UpdateQuizModal />
          <UpdateQuestion />
        </div>
      </FormProvider>
    </>
  );
}
