import { create } from "zustand";

interface QuizQuestion {
  title: string;
  _id: number;
  type: "multiple_choice" | "single_choice" | "true_false";
  marks: number;
  answer_required: boolean;
  details?: string | null | undefined;
  id?: number | null | undefined;
  Media?: any[] | null | undefined;
  answer_explanation?: string | null | undefined;
}

// Define the interface for the store
interface CoursesStore {
  chapterId: number | null;
  setChapterId: (index: number) => void;
  clearChapterId: () => void;

  quizID: number | null;
  questions: QuizQuestion[];
  setQuestions: (questions: QuizQuestion[]) => void;
  appendQuestion: (question: QuizQuestion) => void;
  removeQuestion: (question: QuizQuestion) => void;
  clearQuestions: () => void;

  isNewChapter: boolean;
  isNewLesson: boolean;
  isNewQuiz: boolean;
  isNewAssignment: boolean;
  editItemId: number | null;
  isEditChapter: boolean;
  isNewQuestion: boolean;
  dripSettings: string | null;

  addNewChapter: () => void;
  closeNewChapter: () => void;

  addNewLesson: () => void;
  closeNewLesson: () => void;

  lessonID: number | null;
  isEditLesson: boolean;
  openEditLesson: (id: number) => void;
  closeEditLesson: () => void;

  addNewQuiz: () => void;
  closeNewQuiz: () => void;

  addNewAssignment: () => void;
  closeNewAssignment: () => void;

  assignmentID: number | null;
  isEditAssignment: boolean;
  openEditAssignment: (id: number) => void;
  closeEditAssignment: () => void;

  editQuizID: number | null;
  isEditQuiz: boolean;
  openEditQuiz: (id: number) => void;
  closeEditQuiz: () => void;
  editQuestionID: number | null;
  isEditQuestion: boolean;
  openEditQuestion: (id: number) => void;
  updateQuestion: (updatedQuestion: QuizQuestion) => void;
  closeEditQuestion: () => void;

  setEditItem: (itemId: number | null) => void;

  openEditChapter: () => void;
  closeEditChapter: () => void;

  openNewQuestion: (id: number) => void;
  closeNewQuestion: () => void;

  openDripSettings: (value: string | null) => void;
  closeDripSettings: (value: string | null) => void;
}

// Create the Zustand store
export const useCoursesStore = create<CoursesStore>((set) => ({
  chapterId: null,
  setChapterId: (index) => set({ chapterId: index }),
  clearChapterId: () => set({ chapterId: null }),

  quizID: null,
  questions: [],
  setQuestions: (questions) => set({ questions }),
  appendQuestion: (question) =>
    set((state) => ({
      questions: [...state.questions, question],
    })),
  removeQuestion: (question) =>
    set((state) => ({
      questions: state.questions.filter((q) => q._id !== question._id),
    })),
  clearQuestions: () => set({ questions: [] }),

  isNewChapter: false,
  isNewLesson: false,
  isNewQuiz: false,
  isNewAssignment: false,
  editItemId: null,
  isEditChapter: false,
  isNewQuestion: false,
  dripSettings: null,

  addNewChapter: () => set({ isNewChapter: true }),
  closeNewChapter: () => set({ isNewChapter: false, chapterId: null }),

  addNewLesson: () => set({ isNewLesson: true }),
  closeNewLesson: () => set({ isNewLesson: false, chapterId: null }),

  lessonID: null,
  isEditLesson: false,
  openEditLesson: (id: number) => set({ isEditLesson: true, lessonID: id }),
  closeEditLesson: () =>
    set({ isEditLesson: false, chapterId: null, lessonID: null }),

  addNewQuiz: () => set({ isNewQuiz: true }),
  closeNewQuiz: () => set({ isNewQuiz: false, chapterId: null, questions: [] }),

  addNewAssignment: () => set({ isNewAssignment: true }),
  closeNewAssignment: () => set({ isNewAssignment: false, chapterId: null }),

  assignmentID: null,
  isEditAssignment: false,
  openEditAssignment: (id: number) =>
    set({ isEditAssignment: true, assignmentID: id }),
  closeEditAssignment: () =>
    set({ isEditAssignment: false, chapterId: null, assignmentID: null }),

  editQuizID: null,
  isEditQuiz: false,
  openEditQuiz: (id: number) => set({ isEditQuiz: true, editQuizID: id }),
  closeEditQuiz: () =>
    set({
      isEditQuiz: false,
      chapterId: null,
      editQuizID: null,
      questions: [],
    }),
  editQuestionID: null,
  isEditQuestion: false,
  openEditQuestion: (id: number) =>
    set({ editQuestionID: id, isEditQuestion: true }),
  updateQuestion: (updatedQuestion: QuizQuestion) => {
    set((state) => ({
      questions: state.questions.map((q) =>
        q._id === updatedQuestion._id ? updatedQuestion : q
      ),
    }));
  },
  closeEditQuestion: () => set({ editQuestionID: null, isEditQuestion: false }),

  setEditItem: (itemId) => set({ editItemId: itemId }),
  openEditChapter: () => set({ isEditChapter: true }),
  closeEditChapter: () => set({ isEditChapter: false, chapterId: null }),

  openNewQuestion: (id: number) => set({ isNewQuestion: true, quizID: id }),
  closeNewQuestion: () => set({ isNewQuestion: false, quizID: null }),

  openDripSettings: (type) => set({ dripSettings: type }),
  closeDripSettings: (type) => set({ dripSettings: type }),
}));
