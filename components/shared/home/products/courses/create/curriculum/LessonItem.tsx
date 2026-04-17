"use client";
import TextEditor from "@/components/shared/text-editor/TextEditor";
import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import SelectList from "@/components/ui/SelectList";
import ToggleSwitch from "@/components/ui/ToggleSwitch";
import React, { useEffect, useState } from "react";
import AddRecording from "./AddRecording";
import {
  Controller,
  useFieldArray,
  useForm,
  useFormContext,
} from "react-hook-form";
import {
  CourseLessonSchema,
  TCourseLessonSchema,
  TCourseSchema,
} from "@/schema/course.schema";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import { zodResolver } from "@hookform/resolvers/zod";
import UploadResources from "./UploadResources";
import Schedule from "../details/Schedule";
import LessonSchedule from "../lesson/LessonSchedule";
import { dbTimeToPickerFormat } from "@/lib/helpers";

const mediaOptions = [
  { id: 1, name: "Video", value: "video" },
  { id: 2, name: "Live Session", value: "live_session" },
  { id: 3, name: "Audio", value: "audio" },
  { id: 4, name: "Text", value: "text" },
];

const videoSources = [
  { id: 1, name: "Youtube", value: "youtube" },
  { id: 2, name: "Vimeo", value: "vimeo" },
  { id: 3, name: "Custom Code", value: "custom_code" },
  { id: 4, name: "Upload", value: "upload" },
];

const audioSources = [
  { id: 1, name: "Sound Cloud", value: "sound_cloud" },
  { id: 2, name: "Spotify", value: "spotify" },
  { id: 3, name: "Custom Code", value: "custom_code" },
  { id: 4, name: "Upload", value: "upload" },
];

export default function LessonItem({ isEdit = false }: { isEdit?: boolean }) {
  const {
    chapterId,
    lessonID,
    clearChapterId,
    closeNewLesson,
    closeEditLesson,
  } = useCoursesStore();

  const { control, watch } = useFormContext<TCourseSchema>();

  const chapterIndex = watch("course_chapters", []).findIndex(
    (chapter) => chapter._id === chapterId
  );
  const safeChapterIndex = chapterIndex === -1 ? 0 : chapterIndex;

  const { append, update } = useFieldArray({
    control,
    name: `course_chapters.${safeChapterIndex}.course_lessons`,
    keyName: "uid",
  });

  const formMethods = useForm<TCourseLessonSchema>({
    resolver: zodResolver(CourseLessonSchema),
    defaultValues: {
      _id: Date.now(),
      type: "lesson",
      title: "",
      description: "",
      lesson_type: "video",
      source_type: "youtube",
      source: {
        data: "",
        playback_time: "",
        isFile: false,
      },
      is_scheduled: false,
      schedule_date: null,
      schedule_time: null,
      show_comming_soon: false,
      is_published: false,
      is_public: false,
      resources: null,
    },
  });

  useEffect(() => {
    const watchedLesson =
      watch(`course_chapters.${safeChapterIndex}.course_lessons`) || [];
    if (isEdit && lessonID && watchedLesson?.length > 0) {
      const index = watchedLesson?.findIndex((c) => c._id === lessonID);
      if (index !== -1) {
        const lesson = watchedLesson[index];

        console.log("lesson", lesson);

        formMethods.reset({
          _id: lesson._id,
          id: lesson.id,
          type: "lesson",
          title: lesson.title,
          description: lesson.description,
          is_public: lesson.is_public,
          is_published: lesson.is_published,
          lesson_type: lesson.lesson_type,
          source_type: lesson.source_type,
          is_scheduled: lesson.is_scheduled ? true : false,
          schedule_date: lesson.schedule_date ?? null,
          schedule_time: lesson.schedule_time ?? null,
          show_comming_soon: lesson.show_comming_soon ?? false,
          resources: lesson.resources?.map((r) => {
            // For DB resources, ensure `name` is always present
            if ("isDBImg" in r && r.isDBImg) {
              return {
                ...r,
                name:
                  r.name ||
                  r.title ||
                  r.file_path.split("/").pop() ||
                  "Unknown File",
              };
            }

            // For new uploaded files (File objects)
            if ("name" in r) {
              return r;
            }

            // fallback for any other case
            return {
              course_id: r.course_id,
              id: r.id,
              name: r.title || "Unknown File",
              isDBImg: true,
              url: r.file_path,
              size: r.size || 13495,
              type: r.mine_type,
            };
          }),
          source: {
            data: lesson.source.data,
            playback_time: lesson.source.playback_time,
            isFile: lesson.source.isFile,
          },
        });
      }
    }
  }, [isEdit, lessonID, watch]);

  const handleSave = (data: TCourseLessonSchema) => {
    // console.log(data);

    if (isEdit && lessonID && safeChapterIndex !== -1) {
      const watchedLesson =
        watch(`course_chapters.${safeChapterIndex}.course_lessons`) || [];

      const index = watchedLesson?.findIndex((c) => c._id === lessonID);
      if (index !== -1) {
        const existingResources = (watchedLesson[index].resources || []).map(
          (r) => {
            if (r.isDBImg) {
              return r;
            } else {
              return {
                id: r.id,
                course_id: r.course_id,
                name: r.title,
                url: r.file_path,
                type: r.mine_type,
                isDBImg: true,
                size: 13495,
              };
            }
          }
        );

        const newResources = (data.resources || []).filter((r) => !r.isDBImg);
        const mergedResources = [...existingResources, ...newResources];

        update(index, {
          ...watchedLesson[index],
          title: data.title,
          description: data.description,
          lesson_type: data.lesson_type,
          source_type: data.source_type,
          source: data.source,
          is_published: data.is_published,
          is_public: data.is_public,
          resources: mergedResources,
          is_scheduled: data.is_scheduled,
          schedule_date: data.schedule_date,
          schedule_time: data.schedule_time,
          show_comming_soon: data.show_comming_soon,
        });
      }
      closeEditLesson();
    } else {
      closeNewLesson();
      clearChapterId();
      append({
        _id: Date.now(),
        type: "lesson",
        title: data.title,
        description: data.description,
        lesson_type: data.lesson_type,
        source_type: data.source_type,
        source: data.source,
        is_public: data.is_public,
        is_published: data.is_published,
        resources: data.resources,
        is_scheduled: data.is_scheduled,
        schedule_date: data.schedule_date,
        schedule_time: data.schedule_time,
        show_comming_soon: data.show_comming_soon,
      });
      formMethods.reset();
    }
  };

  const renderPreview = () => {
    const iframeMatch = formMethods
      .watch("source")
      .data.match(/<iframe.*?<\/iframe>/);
    if (iframeMatch) {
      return (
        <div
          className="w-full"
          dangerouslySetInnerHTML={{ __html: iframeMatch[0] }}
        />
      );
    }

    return (
      <p className="text-gray-500 text-sm">
        Preview unavailable. Please add embed code.
      </p>
    );
  };

  return (
    <>
      {/* {JSON.stringify(formMethods.formState.errors, null, 2)} */}
      <div className="z-10 sticky top-0 border-b border-gray-200 bg-white transition-shadow duration-300 mb-5">
        <div className="max-w-7xl mx-auto flex items-center justify-between py-4">
          <p className="font-semibold text-xl px-4">
            {!isEdit ? "New" : "Edit"} Lesson
          </p>
          <Button type="button" onClick={formMethods.handleSubmit(handleSave)}>
            Save
          </Button>
        </div>
      </div>

      <div className="max-w-7xl mx-auto flex items-start gap-5 pb-5">
        <div className="w-full border rounded-xl overflow-hidden">
          <div className="p-4 border-b border-gray-300">
            <p className="font-semibold">Lesson Details</p>
          </div>
          <div className="p-4">
            {/* Title */}
            <div className="mb-5">
              <label
                htmlFor="title"
                className="block mb-1 text-sm font-semibold"
              >
                Title *
              </label>
              <InputField
                className="w-full"
                id="title"
                {...formMethods.register("title")}
                error={formMethods.formState.errors.title?.message}
              />
            </div>

            {/* Description */}
            <div className="mb-5">
              <label
                htmlFor="description"
                className="block mb-1 text-sm font-semibold"
              >
                Description
              </label>
              <Controller
                control={formMethods.control}
                name="description"
                render={({ field }) => (
                  <TextEditor
                    value={field.value || ""}
                    onChange={field.onChange}
                  />
                )}
              />
            </div>

            {/* Conditionally render recording for Live Session */}
            {formMethods.watch("lesson_type") === "live_session" && (
              <div className="mb-5">
                <label className="block mb-1 text-sm font-semibold">
                  Recording
                </label>
                <AddRecording
                  onUploadComplete={(file) => {
                    formMethods.setValue("source.data", file);
                    formMethods.setValue("source.isFile", true);
                  }}
                />
                {formMethods.formState.errors.source?.data?.message && (
                  <p className="text-sm text-red-500 mt-1">
                    {String(formMethods.formState.errors.source?.data?.message)}
                  </p>
                )}
              </div>
            )}

            {/* Video Source */}
            {formMethods.watch("lesson_type") === "video" && (
              <div className="mb-5">
                <label
                  htmlFor="videoSource"
                  className="block mb-1 text-sm font-semibold"
                >
                  Source *
                </label>
                <Controller
                  control={formMethods.control}
                  name="source_type"
                  render={({ field: { value, onChange } }) => {
                    const selected = videoSources.find(
                      (o) => o.value === value
                    );
                    return (
                      <SelectList
                        options={videoSources}
                        value={selected}
                        onChange={(option) => {
                          onChange(option.value);
                          formMethods.setValue("source", {
                            data: "",
                            isFile: false,
                            playback_time: "",
                          });
                        }}
                        className="w-full font-medium"
                      />
                    );
                  }}
                />
              </div>
            )}

            {/* Audio Source */}
            {formMethods.watch("lesson_type") === "audio" && (
              <div className="mb-5">
                <label
                  htmlFor="audioSource"
                  className="block mb-1 text-sm font-semibold"
                >
                  Source *
                </label>
                <Controller
                  control={formMethods.control}
                  name="source_type"
                  render={({ field: { value, onChange } }) => {
                    const selected = audioSources.find(
                      (o) => o.value === value
                    );
                    return (
                      <SelectList
                        options={audioSources}
                        value={selected}
                        onChange={(option) => {
                          onChange(option.value);
                          formMethods.setValue("source", {
                            data: "",
                            isFile: false,
                            playback_time: "",
                          });
                        }}
                        className="w-full font-medium"
                      />
                    );
                  }}
                />
              </div>
            )}

            {/* Video Upload */}
            {formMethods.watch("source_type") === "upload" &&
              formMethods.watch("lesson_type") === "video" && (
                <>
                  <div className="mb-5">
                    <label className="block mb-1 text-sm font-semibold">
                      Recording
                    </label>
                    <AddRecording
                      onUploadComplete={(file) => {
                        formMethods.setValue("source.data", file);
                        formMethods.setValue("source.isFile", true);
                      }}
                    />
                    {formMethods.formState.errors.source?.data?.message && (
                      <p className="text-sm text-red-500 mt-1">
                        {String(
                          formMethods.formState.errors.source?.data?.message
                        )}
                      </p>
                    )}
                  </div>
                  <label className="block mb-1 text-sm font-semibold">
                    Playback Time
                  </label>
                  <InputField
                    placeholder="HH:MM:SS"
                    className="w-full"
                    {...formMethods.register("source.playback_time")}
                    error={
                      formMethods.formState.errors.source?.playback_time
                        ?.message
                    }
                  />
                </>
              )}

            {/* Embed code for Youtube */}
            {formMethods.watch("source_type") === "youtube" &&
              formMethods.watch("lesson_type") === "video" && (
                <>
                  <textarea
                    placeholder='<iframe width="..." height="..." src="..." title="..." frameborder="..." allowfullscreen>...</iframe>'
                    className="w-full border resize-none rounded-md min-h-28 px-3 py-2 text-sm"
                    {...formMethods.register("source.data")}
                  />
                  {formMethods.formState.errors.source?.data?.message && (
                    <p className="text-sm text-red-500 mt-1">
                      {String(
                        formMethods.formState.errors.source?.data?.message
                      )}
                    </p>
                  )}
                  <div className="w-full border rounded-md p-5 mt-5 flex items-center justify-center min-h-44 mb-5">
                    {renderPreview()}
                  </div>
                  <label className="block mb-1 text-sm font-semibold">
                    Playback Time
                  </label>
                  <InputField
                    placeholder="HH:MM:SS"
                    className="w-full"
                    {...formMethods.register("source.playback_time")}
                    error={
                      formMethods.formState.errors.source?.playback_time
                        ?.message
                    }
                  />
                </>
              )}
          </div>
        </div>

        <div className="w-96 min-w-96">
          <div className="border rounded-xl">
            <div className="p-4 border-b border-gray-300">
              <p className="font-semibold">Lesson Settings</p>
            </div>
            <div className="p-4">
              {/* Publish toggle */}
              {!formMethods.watch("is_scheduled") && (
                <>
                  <div className="flex items-center justify-between">
                    <label className="text-sm font-semibold">Publish</label>
                    <Controller
                      control={formMethods.control}
                      name="is_published"
                      render={({ field: { value, onChange } }) => (
                        <ToggleSwitch
                          checked={value}
                          onChange={(val) => onChange(val)}
                        />
                      )}
                    />
                  </div>
                  <p className="text-sm text-gray-600 font-medium">
                    Make lesson as publish
                  </p>
                </>
              )}

              {/* Public toggle */}
              <div className="flex items-center justify-between mt-3">
                <label className="text-sm font-semibold">Public</label>
                <Controller
                  control={formMethods.control}
                  name="is_public"
                  render={({ field: { value, onChange } }) => (
                    <ToggleSwitch
                      checked={value}
                      onChange={(val) => onChange(val)}
                    />
                  )}
                />
              </div>
              <p className="text-sm text-gray-600 font-medium">
                Make lesson public for students
              </p>

              {/* Lesson Type Select */}
              <div className="my-4">
                <label className="text-sm font-semibold block">
                  Lesson Type <span className="text-red-500">*</span>
                </label>
                <Controller
                  control={formMethods.control}
                  name="lesson_type"
                  render={({ field: { value, onChange } }) => {
                    const selected = mediaOptions.find(
                      (o) => o.value === value
                    );
                    return (
                      <SelectList
                        options={mediaOptions}
                        value={selected}
                        onChange={(option) => {
                          onChange(option.value);
                          formMethods.setValue("source", {
                            data: "",
                            isFile: false,
                            playback_time: "",
                          });
                        }}
                        className="w-full font-medium"
                      />
                    );
                  }}
                />
              </div>
            </div>
          </div>

          <div className="mt-4">
            <LessonSchedule formMethods={formMethods} />
          </div>

          <div className="border rounded-xl mt-5">
            <div className="p-4 border-b border-gray-300">
              <p className="font-semibold">Lesson Resources</p>
            </div>
            <div className="p-4">
              <Controller
                control={formMethods.control}
                name="resources"
                render={({ field: { value, onChange } }) => (
                  <UploadResources
                    value={value ?? []}
                    onFilesSelected={(files) => {
                      onChange(files);
                      formMethods.trigger("resources");
                    }}
                  />
                )}
              />
              {formMethods.formState.errors.resources?.message && (
                <p className="text-sm text-red-500 mt-1">
                  {String(formMethods.formState.errors.resources?.message)}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
