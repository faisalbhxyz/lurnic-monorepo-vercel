import { z } from "zod";

export const QuizQuestionSchema = z.object({
  _id: z.coerce.number(),
  id: z.coerce.number().optional().nullable(),
  title: z
    .string({ required_error: "Title is required" })
    .trim()
    .min(1, {
      message: "Title is required",
    })
    .max(200, { message: "Title should not exceed 200 characters" }),
  details: z.string().optional().nullable(),
  media: z.array(z.any()).optional().nullable(),
  type: z.enum(["multiple_choice", "single_choice", "true_false"]),
  marks: z.coerce.number({ required_error: "Marks is required" }).gte(0),
  answer_required: z.boolean(),
  answer_explanation: z.string().trim().optional().nullable(),
});

export type TQuizQuestionSchema = z.infer<typeof QuizQuestionSchema>;

export const CourseQuizSchema = z.object({
  _id: z.coerce.number(),
  id: z.coerce.number().optional().nullable(),
  type: z.literal("quiz"),
  title: z
    .string({ required_error: "Title is required" })
    .trim()
    .min(1, {
      message: "Title is required",
    })
    .max(200, { message: "Title should not exceed 200 characters" }),
  instructions: z
    .string({ required_error: "Instructions is required" })
    .trim()
    .min(1, {
      message: "Instructions is required",
    }),
  is_published: z.boolean(),
  randomize_questions: z.boolean(),
  single_quiz_view: z.boolean(),
  time_limit: z.coerce
    .number({ required_error: "Time limit is required" })
    .gte(0),
  time_limit_option: z.enum(["minutes", "hours", "days", "weeks", "months"]),
  total_visible_questions: z.coerce
    .number({ required_error: "Total visible questions is required" })
    .gte(0),
  reveal_answers: z.boolean(),
  enable_retry: z.boolean(),
  retry_attempts: z.coerce
    .number({ required_error: "Retry attempts is required" })
    .gte(0),
  minimum_pass_percentage: z.coerce
    .number({ required_error: "Minimum pass percentage is required" })
    .gte(0),
  questions: z.array(QuizQuestionSchema),
});

export type TCourseQuizSchema = z.infer<typeof CourseQuizSchema>;

export const CourseAssignmentSchema = z.object({
  _id: z.coerce.number(),
  id: z.coerce.number().optional().nullable(),
  type: z.literal("assignment"),
  title: z
    .string({ required_error: "Title is required" })
    .trim()
    .min(1, {
      message: "Title is required",
    })
    .max(200, { message: "Title should not exceed 200 characters" }),
  instructions: z
    .string({ required_error: "Instructions is required" })
    .trim()
    .min(1, {
      message: "Instructions is required",
    }),
  attachments: z.array(z.any()).optional().nullable(),
  is_published: z.boolean(),
  time_limit: z.coerce
    .number({ required_error: "Time limit is required" })
    .gte(0),
  time_limit_option: z.enum(["minutes", "hours", "days", "weeks", "months"]),
  file_upload_limit: z.coerce
    .number({ required_error: "File upload limit is required" })
    .gte(0),
  total_marks: z.coerce
    .number({ required_error: "Total marks is required" })
    .gte(0),
  minimum_pass_marks: z.coerce
    .number({ required_error: "Minimum pass marks is required" })
    .gte(0),
});

export type TCourseAssignmentSchema = z.infer<typeof CourseAssignmentSchema>;

export const CourseLessonSchema = z
  .object({
    _id: z.coerce.number(),
    id: z.coerce.number().optional().nullable(),
    type: z.literal("lesson"),
    title: z
      .string({ required_error: "Title is required" })
      .min(1, { message: "Title is required" })
      .max(100, { message: "Title should not exceed 100 characters" }),
    description: z.string().optional().nullable(),
    lesson_type: z.enum(["video", "live_session", "audio", "text"]),
    source_type: z.enum([
      "youtube",
      "vimeo",
      "sound_cloud",
      "spotify",
      "recording",
      "custom_code",
      "upload",
    ]),
    source: z.object({
      data: z.any(),
      playback_time: z.string().optional().nullable(),
      isFile: z.boolean(),
    }),
    is_published: z.boolean(),
    is_public: z.boolean(),
    is_scheduled: z.boolean(),
    schedule_date: z.string().optional().nullable(),
    schedule_time: z.string().optional().nullable(),
    show_comming_soon: z.boolean(),
    resources: z
      .array(
        z.any().refine((file) => {
          if (!file) return true; // Allow empty
          if (file.isDBImg) return true;
          if (!(file instanceof File)) return false;
          return file.size <= 2 * 1024 * 1024;
        }, "Max image size is 2MB.")
      )
      .optional()
      .nullable(),
  })
  .superRefine(({ source }, ctx) => {
    if (source.isFile) {
      if (!source.data) {
        ctx.addIssue({
          code: "custom",
          path: ["source", "data"],
          message: "A file must be provided.",
        });
      }
    } else {
      if (
        source.data === null ||
        source.data === undefined ||
        (typeof source.data === "string" && source.data.trim() === "")
      ) {
        ctx.addIssue({
          code: "custom",
          path: ["source", "data"],
          message: "This field is required.",
        });
      }
    }
  });

export type TCourseLessonSchema = z.infer<typeof CourseLessonSchema>;

export const CourseChapterSchema = z.object({
  _id: z.coerce.number(),
  id: z.coerce.number().optional().nullable(),
  position: z.coerce.number().gte(0),
  title: z
    .string({ required_error: "Title is required" })
    .min(1, {
      message: "Title is required",
    })
    .max(100, { message: "Title should not exceed 100 characters" }),
  description: z.string().optional().nullable(),
  access: z.enum(["draft", "published"]),
  course_lessons: z.array(CourseLessonSchema).nullable().optional(),
  quizzes: z.array(CourseQuizSchema).optional().nullable(),
  assignments: z.array(CourseAssignmentSchema).optional().nullable(),
});

export type TCourseChapterSchema = z.infer<typeof CourseChapterSchema>;

const GeneralSettingsSchema = z.object({
  difficulty_level: z.enum(["all", "beginner", "intermediate", "expert"]),
  maximum_student: z.coerce
    .number({
      invalid_type_error: "Maximum student is required",
      required_error: "Maximum student is required",
      message: "Maximum student is required",
    })
    .gte(0),
  language: z.string().nullable().optional(),
  category_id: z.coerce
    .number({
      required_error: "Please select a category",
      invalid_type_error: "Please select a category",
    })
    .gt(0, {
      message: "Please select a category",
    }),
  sub_category_id: z.coerce
    .number({
      invalid_type_error: "Please select a sub category",
    })
    .gt(0, {
      message: "Please select a sub category",
    })
    .optional()
    .nullable(),
  duration: z.string().nullable().optional(),
});

export const CourseSchema = z.object({
  title: z
    .string({ required_error: "Title is required" })
    .min(1, {
      message: "Title is required",
    })
    .max(100, { message: "Title should not exceed 100 characters" }),
  summary: z
    .string({ required_error: "Summary is required" })
    .min(1, { message: "Summary is required" })
    .refine((val) => val.trim().split(/\s+/).length <= 50, {
      message: "Summary should not exceed 50 words",
    }),
  description: z.string().nullish(),
  visibility: z.enum(["public", "private", "protected"]),
  is_scheduled: z.boolean(),
  schedule_date: z.string().optional().nullable(),
  schedule_time: z.string().optional().nullable(),
  featured_image: z
    .any()
    .refine((file) => {
      if (!file) return true; // Allow empty
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;
      return file.size <= 2 * 1024 * 1024;
    }, "Max image size is 2MB.")
    .refine((file) => {
      if (!file) return true; // Allow empty
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;
      return [
        "image/png",
        "image/jpg",
        "image/jpeg",
        "image/webp",
        "image/gif",
      ].includes(file.type); // Check file type
    }, "Only .png, .jpg & .jpeg formats are supported.")
    .refine((file) => {
      if (!file) return true;
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;

      return new Promise<boolean>((resolve) => {
        const img = document.createElement("img");
        img.src = URL.createObjectURL(file);
        img.onload = () => {
          const isValid = img.width <= 1920 && img.height <= 1080;
          resolve(isValid);
          URL.revokeObjectURL(img.src); // cleanup
        };
        img.onerror = () => {
          resolve(false);
          URL.revokeObjectURL(img.src); // cleanup
        };
      });
    }, "Image must be 1920x1080 pixels or smaller."),
  intro_video: z
    .object({
      type: z.enum(["youtube", "vimeo"]),
      source: z
        .string()
        .trim()
        .url({ message: "Invalid URL" })
        .optional()
        .nullable(),
    })
    .optional()
    .nullable(),
  pricing_model: z.enum(["free", "paid"]),
  regular_price: z.coerce.number().gte(0),
  sale_price: z.coerce.number().gte(0),
  show_comming_soon: z.boolean(),
  tags: z.array(z.string()).optional().nullable(),
  overview: z.array(z.string()).optional().nullable(),
  author_id: z.coerce.number(),
  course_chapters: z
    .array(CourseChapterSchema)
    .min(1, { message: "At least one chapter is required" }),
  course_instructors: z
    .array(z.coerce.number(), {
      message: "At least one instructor is required",
    })
    .min(1, {
      message: "At least one instructor is required",
    }),
  general_settings: GeneralSettingsSchema,
});

export type TCourseSchema = z.infer<typeof CourseSchema>;
