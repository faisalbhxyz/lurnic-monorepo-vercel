import {
  TCourseAssignmentSchema,
  TCourseChapterSchema,
  TCourseLessonSchema,
  TCourseQuizSchema,
} from "@/schema/course.schema";

export type ChapterItem =
  | TCourseLessonSchema
  | TCourseQuizSchema
  | TCourseAssignmentSchema;

const TYPE_ORDER: Record<ChapterItem["type"], number> = {
  lesson: 0,
  quiz: 1,
  assignment: 2,
};

export function getItemDragId(item: ChapterItem): string {
  return `${item.type}:${item._id}`;
}

export function buildChapterItems(chapter: TCourseChapterSchema): ChapterItem[] {
  const items: ChapterItem[] = [
    ...(chapter.course_lessons ?? []),
    ...(chapter.quizzes ?? []),
    ...(chapter.assignments ?? []),
  ];

  return items.sort((a, b) => {
    const posA = a.position ?? Number.MAX_SAFE_INTEGER;
    const posB = b.position ?? Number.MAX_SAFE_INTEGER;
    if (posA !== posB) return posA - posB;
    return TYPE_ORDER[a.type] - TYPE_ORDER[b.type];
  });
}

export function splitChapterItems(items: ChapterItem[]) {
  return {
    course_lessons: items.filter(
      (item): item is TCourseLessonSchema => item.type === "lesson"
    ),
    quizzes: items.filter(
      (item): item is TCourseQuizSchema => item.type === "quiz"
    ),
    assignments: items.filter(
      (item): item is TCourseAssignmentSchema => item.type === "assignment"
    ),
  };
}

export function withChapterItemPositions(items: ChapterItem[]): ChapterItem[] {
  return items.map((item, index) => ({ ...item, position: index }));
}

/** Assign sequential positions using legacy display order when saved positions are missing. */
export function normalizeChapterItemPositions(
  chapter: TCourseChapterSchema
): TCourseChapterSchema {
  const items: ChapterItem[] = [
    ...(chapter.course_lessons ?? []),
    ...(chapter.quizzes ?? []),
    ...(chapter.assignments ?? []),
  ];

  if (items.length === 0) {
    return chapter;
  }

  const positions = items.map((item) => item.position);
  const usesSavedOrder =
    positions.every((position) => position != null) &&
    new Set(positions).size === items.length;

  if (usesSavedOrder) {
    return {
      ...chapter,
      ...splitChapterItems(buildChapterItems(chapter)),
    };
  }

  const lessons = [...(chapter.course_lessons ?? [])].sort(
    (a, b) => (a.position ?? 0) - (b.position ?? 0)
  );
  const legacyItems: ChapterItem[] = [
    ...lessons,
    ...(chapter.quizzes ?? []),
    ...(chapter.assignments ?? []),
  ];

  return {
    ...chapter,
    ...splitChapterItems(withChapterItemPositions(legacyItems)),
  };
}

export function nextChapterItemPosition(
  chapter: Pick<
    TCourseChapterSchema,
    "course_lessons" | "quizzes" | "assignments"
  >
): number {
  return buildChapterItems(chapter).length;
}
