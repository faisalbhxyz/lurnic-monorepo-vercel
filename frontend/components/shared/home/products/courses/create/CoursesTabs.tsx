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
import { dbTimeToPickerFormat } from "@/lib/helpers";

const tabs = [
  { id: 1, label: "Details" },
  { id: 2, label: "Curriculum" },
  // { id: 3, label: "Quizzes" },
  // { id: 4, label: "Assignments" },
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
        },
      ],
      general_settings: {
        difficulty_level: "all",
        maximum_student: 0,
        language: "english",
        duration: null,
        sub_category_id: null,
      },
    },
  });

  useEffect(() => {
    if (isEdit && courseDetails) {
      console.log("courseDetails", courseDetails);

      const instructors = courseDetails.course_instructors.map(
        (ins) => ins.instructor.id
      );

      const chapters = courseDetails.course_chapters.map((chapter) => {
        return {
          ...chapter,
          id: chapter.id,
          _id: chapter.id,
          course_lessons: chapter.course_lessons.map((lesson) => {
            const source = lesson.source || {};
            return {
              ...lesson,
              id: lesson.id,
              _id: lesson.id,
              type: "lesson",
              source: {
                data: source.data.data,
                playback_time: source.data.playback_times,
                isFile: false,
              },
            };
          }),
          assignments: chapter.assignments.map((assignment) => {
            return {
              ...assignment,
              id: assignment.id,
              _id: assignment.id,
              type: "assignment",
            };
          }),
          quizzes: chapter.quizzes.map((quiz) => {
            return {
              ...quiz,
              id: quiz.id,
              _id: quiz.id,
              type: "quiz",
              questions: quiz.questions.map((q) => {
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
        course_chapters: chapters.map((chapter) => ({
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
        })) as any,
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
          sub_category_id: courseDetails.general_settings.sub_category_id,
        },
      });
    }
  }, [isEdit, courseDetails]);

  const renderContent = () => {
    switch (activeTab) {
      case "Details":
        return <Basics />;
      case "Curriculum":
        return <Curriculum />;
      // case "Quizzes":
      //   return (
      //     <div className="border p-5 rounded-lg mt-5">
      //       <QuizTable />
      //     </div>
      //   );
      // case "Assignments":
      //   return (
      //     <div className="border p-5 rounded-lg mt-5">
      //       <AssignmentTable />
      //     </div>
      //   );
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

  const handleSave = (data: TCourseSchema) => {
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
    const fd = new FormData();
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
    fd.append("course_chapters", JSON.stringify(data.course_chapters));
    fd.append("general_settings", JSON.stringify(data.general_settings));
    fd.append("course_instructors", JSON.stringify(data.course_instructors));
    data.course_chapters.forEach((chapter, chapterIndex) => {
      chapter.course_lessons?.forEach((lesson, lessonIndex) => {
        lesson.resources?.forEach((file, resourceIndex) => {
          if (!file.isDBImg) {
            fd.append(
              `resources[${chapterIndex}][${lessonIndex}][]`,
              file,
              file.name
            );
          }
        });
      });
    });

    if (isEdit && courseDetails) {
      // console.log("[EDIT]", data);
      axiosInstance
        .put(`/private/course/update/${courseDetails.id}`, fd, {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: `Bearer ${session?.accessToken}`,
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
          toast.error(error.response.data.error || "Something went wrong.");
        })
        .finally(() => setLoading(false));
    } else {
      axiosInstance
        .post("/private/course/create", fd, {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: `Bearer ${session?.accessToken}`,
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
          toast.error(error.response.data.error || "Something went wrong.");
        })
        .finally(() => setLoading(false));
    }
  };

  return (
    <>
      <div className="border-b border-gray-300 flex space-x-4 mb-4">
        {tabs.map((tab) => (
          <button
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
        <form onSubmit={formMethods.handleSubmit(handleSave)}>
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
                type="submit"
                disabled={loading}
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
        </form>
      </FormProvider>
    </>
  );
}
