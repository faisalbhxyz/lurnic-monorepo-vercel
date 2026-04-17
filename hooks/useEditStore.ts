import { create } from "zustand";

interface EditStore {
  editID: number | null;
  setEditID: (id: number) => void;

  // for users
  isUserEditOpen: boolean;
  toggleUserEdit: (state: boolean) => void;

  // for students
  isStudentEditOpen: boolean;
  toggleStudentEdit: (state: boolean) => void;

  // for instructor
  isInstructorEditOpen: boolean;
  toggleInstructorEdit: (state: boolean) => void;

  // for payment methods
  refreshPaymentMethodsCount: number;
  refreshPaymentMethods: () => void;
  isPaymentMethodEditOpen: boolean;
  togglePaymentMethodEdit: (state: boolean) => void;
}

export const useEditStore = create<EditStore>((set) => ({
  editID: null,
  setEditID: (id) => set({ editID: id }),

  // for users
  isUserEditOpen: false,
  toggleUserEdit: (state) =>
    set((prev) => ({
      isUserEditOpen: state,
      editID: state ? prev.editID : null,
    })),

  // for students
  isStudentEditOpen: false,
  toggleStudentEdit: (state) =>
    set((prev) => ({
      isStudentEditOpen: state,
      editID: state ? prev.editID : null,
    })),

  // for instructor
  isInstructorEditOpen: false,
  toggleInstructorEdit: (state) =>
    set((prev) => ({
      isInstructorEditOpen: state,
      editID: state ? prev.editID : null,
    })),

  // for payment methods
  refreshPaymentMethodsCount: 0,
  refreshPaymentMethods: () =>
    set((prev) => ({
      refreshPaymentMethodsCount: prev.refreshPaymentMethodsCount + 1,
    })),
  isPaymentMethodEditOpen: false,
  togglePaymentMethodEdit: (state) =>
    set((prev) => ({
      isPaymentMethodEditOpen: state,
      editID: state ? prev.editID : null,
    })),
}));
