interface ICategory {
  id: number;
  name: string;
  slug: string;
  description: string | null;
  thumbnail: string | null;
  created_at: string;
  updated_at: string;
  subcategories: ISubCategory[];
}
interface ISubCategory {
  id: number;
  name: string;
  slug: string;
  description: string | null;
  thumbnail: string | null;
  created_at: string;
  updated_at: string;
  category_id: number;
  category: ICategory;
}

interface IRole {
  id: number;
  name: string;
  permissions: string[] | null;
  created_at: string;
  updated_at: string;
}

interface IUser {
  id: number;
  user_id: string;
  name: string;
  phone?: string | null;
  email: string;
  status: string;
  role_id: number;
  role: IRole;
  created_at: string;
  updated_at: string;
}

interface IInstructor {
  id: number;
  user_id: string;
  first_name: string;
  last_name?: string | null;
  phone?: string | null;
  email: string;
  status: string;
  created_at: string;
  updated_at: string;
}

interface IStudent {
  id: number;
  user_id: string;
  first_name: string;
  last_name?: string | null;
  phone?: string | null;
  email: string;
  status: string;
  created_at: string;
  updated_at: string;
  enrollments?: IEnrollment[] | null;
}

interface ICourseLesson {
  _id: number;
  type: "lesson";
  title: string;
  description?: string | null;
  lesson_type: "video" | "live_session" | "audio" | "text";
  source_type:
    | "youtube"
    | "vimeo"
    | "sound_cloud"
    | "spotify"
    | "recording"
    | "custom_code"
    | "upload";
  source: {
    data: File | string;
    playback_time?: string | null;
    isFile: boolean;
  };
  is_published: boolean;
  is_public: boolean;
  resources?: File[] | null;
  created_at: string;
  updated_at: string;
}

interface ICourseChapter {
  id: number;
  position: number;
  title: string;
  description?: string | null;
  access: "draft" | "published";
  course_lessons?: ICourseLesson[] | null;
  created_at: string;
  updated_at: string;
}

interface IGeneralSettings {
  difficulty_level: "all" | "beginner" | "intermediate" | "expert";
  maximum_student: number;
  language?: string | null;
  category_id: number;
  duration?: string | null;
  created_at: string;
  updated_at: string;
}

interface ICourse {
  title: string;
  summary: string;
  description?: string | null;
  visibility: "public" | "private" | "protected";
  is_scheduled: boolean;
  schedule_date?: Date | null;
  schedule_time?: Date | null;
  featured_image?: string | null;
  intro_video?: string | null;
  pricing_model: "free" | "paid";
  regular_price: number;
  sale_price: number;
  show_comming_soon: boolean;
  tags?: string[] | null;
  overview?: string[] | null;
  author_id: number;
  course_chapters: CourseChapter[];
  course_instructors: Instructor[];
  general_settings: GeneralSettings;
  created_at: string;
  updated_at: string;
}

type Visibility = "public" | "private" | "protected";

type PricingModel = "free" | "paid";

type DifficultyLevel = "all" | "beginner" | "intermediate" | "expert";

type Access = "draft" | "published";

type LessonType = "video" | "live_session" | "audio" | "text";

type LessonSourceType =
  | "youtube"
  | "vimeo"
  | "custom_code"
  | "upload"
  | "sound_cloud"
  | "spotify";

interface IntroVideo {
  type: string;
  source: string;
}

interface Source {
  data: string;
  is_file: boolean;
  playback_times?: string | null;
}

interface CourseDetails {
  id: number;
  title: string;
  summary: string;
  description?: string | null;
  visibility: Visibility;
  is_scheduled?: boolean | null;
  schedule_date?: string | null; // ISO date string
  schedule_time?: string | null; // ISO time string or string format
  show_comming_soon?: boolean | null;
  featured_image?: string | null;
  intro_video?: IntroVideo | null;
  pricing_model: PricingModel;
  regular_price?: number | null;
  sale_price?: number | null;
  tags: any; // JSON stored as any, adjust if you know the shape
  overview: any; // JSON stored as any, adjust if you know the shape
  author_id: number;
  author: IUser; // You will need to define this interface separately
  created_at: string; // ISO datetime string
  updated_at: string; // ISO datetime string
  tenant_id?: number;
  tenant?: Tenant; // Define Tenant interface separately if needed
  course_chapters: CourseChapter[];
  general_settings: CourseGeneralSettings;
  course_instructors: CourseInstructor[];
}

interface CourseChapter {
  id: number;
  position: number;
  title: string;
  description?: string | null;
  access: Access;
  created_at: string;
  updated_at: string;
  course_id: number;
  course_lessons: CourseLesson[];
  assignments: CourseAssignment[];
  quizzes: CourseQuiz[];
}

interface CourseLesson {
  id: number;
  title: string;
  description?: string | null;
  lesson_type: LessonType;
  source_type: LessonSourceType;
  source: { data: Source };
  is_published: boolean;
  is_public: boolean;
  resources?: Record<string, string> | null; // filename, mimetype, url, size
  position: number;
  is_scheduled?: boolean | null;
  schedule_date?: string | null; // ISO date string
  schedule_time?: string | null; // ISO time string or string format
  show_comming_soon?: boolean | null;
  created_at: string;
  updated_at: string;
  chapter_id: number;
}
interface CourseAssignment {
  id: number;
  course_id: number;
  chapter_id: number;
  title: string;
  instructions: string;
  attachments: any | null;
  is_published: boolean;
  total_marks: number;
  minimum_pass_marks: number;
  time_limit: number;
  time_limit_option: "minutes" | "hours" | "days" | "weeks" | "months";
  file_upload_limit: number;
  created_at: string;
  updated_at: string;
}
interface CourseQuiz {
  id: number;
  title: string;
  instructions: string;
  minimum_pass_percentage: number;
  enable_retry: boolean;
  retry_attempts: number;
  randomize_questions: boolean;
  reveal_answers: boolean;
  single_quiz_view: boolean;
  time_limit: number;
  time_limit_option: string;
  total_visible_questions: number;
  is_published: boolean;
  chapter_id: number;
  course_id: number;
  created_at: string;
  updated_at: string;
  questions: QuizQuestion[];
}

interface QuizQuestion {
  id: number;
  quiz_id: number;
  title: string;
  details: string;
  type: string;
  marks: number;
  answer_explanation: string | null;
  answer_required: boolean;
  media: any[];
  created_at: string;
  updated_at: string;
}

interface CourseGeneralSettings {
  id: number;
  course_id: number;
  difficulty_level?: DifficultyLevel | null;
  maximum_student?: number | null;
  language?: string | null;
  category_id: number;
  sub_category_id: number;
  category: ICategory; // Define Category interface separately
  duration?: string | null;
  created_at: string;
  updated_at: string;
}

interface CourseInstructor {
  id: number;
  course_id: number;
  instructor_id: number;
  instructor: IInstructor; // Define Instructor interface separately
  created_at: string;
  updated_at: string;
}

interface Enrollment {
  id: number;
  student_id: number;
  student: Pick<IStudent, "id" | "first_name" | "last_name" | "email">;
  course_id: number;
  course: Pick<CourseDetails, "id" | "title">;
  created_at: string;
  updated_at: string;
}

interface IBanner {
  id: number;
  title?: string | null;
  url?: string | null;
  image: string;
  created_at: string;
  updated_at: string;
}

interface IOrder {
  id: number;
  student_id: number;
  student: IStudent;
  course_id: number;
  course: CourseDetails;
  discount_type: string;
  discount: number;
  total: number;
  payment_status: "paid" | "unpaid";
  invoice_id: string;
  payment_type: string;
  admin_note?: string | null;
  customer_note?: string | null;
  payment_method: string | null;
  transaction_id: string | null;
  created_at: string;
  updated_at: string;
}

interface GeneralSettings {
  id: number;
  org_name: string;
  logo: string | null;
  favicon: string | null;
  student_prefix: string;
  teacher_prefix: string;
  created_at: string;
  updated_at: string;
}

interface IPaymentMethods {
  id: number;
  title: string;
  image: string | null;
  instruction: string;
  status: boolean;
  created_at: string;
  updated_at: string;
}
