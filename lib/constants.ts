export const pageFeatures = [
  {
    id: 1,
    name: "Instructor Info",
    enable: true,
  },
  {
    id: 2,
    name: "Q&A",
    enable: false,
  },
  {
    id: 3,
    name: "Author",
    enable: false,
  },
  {
    id: 4,
    name: "Level",
    enable: false,
  },
  {
    id: 5,
    name: "Social Share",
    enable: false,
  },
  {
    id: 6,
    name: "Duration",
    enable: false,
  },
  {
    id: 7,
    name: "Total Enrolled",
    enable: false,
  },
  {
    id: 8,
    name: "Update Date",
    enable: false,
  },
  {
    id: 9,
    name: "Progress Bar",
    enable: false,
  },
  {
    id: 10,
    name: "Material",
    enable: false,
  },
  {
    id: 11,
    name: "About",
    enable: false,
  },
  {
    id: 12,
    name: "Description",
    enable: false,
  },
  {
    id: 13,
    name: "Benefits",
    enable: false,
  },
  {
    id: 14,
    name: "Requirements",
    enable: false,
  },
  {
    id: 15,
    name: "Target Audience",
    enable: false,
  },
  {
    id: 16,
    name: "Announcements",
    enable: false,
  },
  {
    id: 17,
    name: "Review",
    enable: false,
  },
];

export const presetColors = [
  {
    name: "Default",
    colors: {
      primaryColor: "#155dfc",
      hoverColor: "#1447e6",
      textColor: "#000000",
      gray: "#99a1af",
    },
  },
  {
    name: "Landscape",
    colors: {
      primaryColor: "#00bba7",
      hoverColor: "#009689",
      textColor: "#000000",
      gray: "#99a1af",
    },
  },
  {
    name: "Ocean",
    colors: {
      primaryColor: "#8e51ff",
      hoverColor: "#7008e7",
      textColor: "#000000",
      gray: "#99a1af",
    },
  },
  {
    name: "Custom",
    colors: {
      primaryColor: "#e12afb",
      hoverColor: "#a800b7",
      textColor: "#000000",
      gray: "#99a1af",
    },
  },
];

export const items = [
  { id: 2, type: "Lesson", name: "New Lesson", access: "Draft" },
  { id: 3, type: "Quiz", name: "New Quiz", access: "Published" },
  { id: 4, type: "Assignment", name: "New Assignment", access: "Draft" },
];

export const initialChapters = [
  {
    id: 1,
    name: "New Chapter 1",
    access: "Published",
    items,
  },
  { id: 2, name: "New Chapter 2", access: "Draft", items: [] },
  { id: 3, name: "New Chapter 3", access: "Published", items: [] },
  { id: 4, name: "New Chapter 4", access: "Draft", items: [] },
];
