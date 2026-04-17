"use client";

import React, { useEffect, useState } from "react";

// DnD Kit
import {
  closestCenter,
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";

// Icons
import {
  MdChecklist,
  MdDragIndicator,
  MdOutlineFileCopy,
} from "react-icons/md";
import { FiBookOpen } from "react-icons/fi";

// Components & Hooks
import Button from "@/components/ui/Button";
import AddItem from "./AddItem";
import ChapterActions from "./ChapterActions";
import ItemActions from "./ItemActions";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import { LuPlus } from "react-icons/lu";
import { useFormContext } from "react-hook-form";
import { CourseChapterSchema, TCourseSchema } from "@/schema/course.schema";
import { z } from "zod";

// Main Curriculum Component
export default function Curriculum() {
  const { addNewChapter } = useCoursesStore();

  const { watch } = useFormContext<TCourseSchema>();

  const chapters = watch("course_chapters", []);

  const [activeId, setActiveId] = useState<number | null>(null);
  const [isClient, setIsClient] = useState(false);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: { distance: 5 },
      eventOptions: { passive: false },
    })
  );

  useEffect(() => {
    setIsClient(true);
  }, []);

  const handleChapterDragStart = (event: DragStartEvent) => {
    console.log("event", event);
    setActiveId(event.active.id as number);
  };

  const handleChapterDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveId(null);

    if (!over || active.id === over.id) return;

    const oldIndex = chapters.findIndex((ch) => ch._id === active.id);
    const newIndex = chapters.findIndex((ch) => ch._id === over.id);

    // if (oldIndex !== -1 && newIndex !== -1) {
    //   setChapters((prev) => arrayMove(prev, oldIndex, newIndex));
    // }
  };

  const activeChapter = chapters.find((ch) => ch._id === activeId);

  const handleItemReorder = (chapterId: number, newItems: ICourseLesson[]) => {
    // setChapters((prev) =>
    //   prev.map((ch) => (ch.id === chapterId ? { ...ch, items: newItems } : ch))
    // );
  };

  if (!isClient) return null;

  return (
    <div className="w-full bg-gray-100 mx-auto border mt-5 rounded-xl p-5">
      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragStart={handleChapterDragStart}
        onDragEnd={handleChapterDragEnd}
      >
        <SortableContext
          items={chapters.map((ch) => ch._id)}
          strategy={verticalListSortingStrategy}
        >
          <div className="mb-8 space-y-4">
            {chapters.map((chapter) => (
              <SortableChapter
                key={chapter._id}
                chapter={chapter}
                onItemReorder={handleItemReorder}
              />
            ))}
          </div>
        </SortableContext>

        <DragOverlay>
          {activeChapter && (
            <SortableChapter chapter={activeChapter} onItemReorder={() => {}} />
          )}
        </DragOverlay>
      </DndContext>

      <div className="flex flex-col items-center justify-center mt-6">
        <Button onClick={addNewChapter} className="p-1.5">
          <LuPlus size={22} />
        </Button>
        <span className="text-sm font-medium text-gray-700 mt-2">
          Add Chapter
        </span>
      </div>
    </div>
  );
}

// Sortable Chapter Component
function SortableChapter({
  chapter,
  onItemReorder,
}: {
  chapter: z.infer<typeof CourseChapterSchema>;
  onItemReorder: (chapterId: number, items: ICourseLesson[]) => void;
}) {
  // const [chapterItems, setChapterItems] = useState<ItemType[]>(chapter.items);
  const [activeId, setActiveId] = useState<number | null>(null);

  const chapterItems = [
    ...(chapter.course_lessons ?? []),
    ...(chapter.quizzes ?? []),
    ...(chapter.assignments ?? []),
  ];

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: { distance: 5 },
      eventOptions: { passive: false },
    })
  );

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: chapter._id });

  const style: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const handleItemDragStart = (event: DragStartEvent) => {
    setActiveId(event.active.id as number);
  };

  const handleItemDragEnd = (event: DragEndEvent) => {
    // const { active, over } = event;
    // setActiveId(null);
    // if (!over || active.id === over.id) return;
    // const oldIndex = chapterItems.findIndex((item) => item.id === active.id);
    // const newIndex = chapterItems.findIndex((item) => item.id === over.id);
    // if (oldIndex !== -1 && newIndex !== -1) {
    //   onItemReorder(chapter.id, arrayMove(chapterItems, oldIndex, newIndex));
    // }
  };

  const activeItem =
    chapter.course_lessons &&
    chapter.course_lessons.find((item) => item._id === activeId);

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      className="flex items-center gap-1"
    >
      <div {...listeners} className="cursor-grab text-gray-500">
        <MdDragIndicator />
      </div>
      <div className="w-full bg-white border rounded-xl p-3">
        {/* Chapter Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="font-medium">{chapter.title}</span>
            <span className="ml-5 text-xs bg-gray-200 px-2 py-0.5 rounded-md">
              {chapter.access}
            </span>
          </div>
          <div className="flex items-center">
            <AddItem id={chapter._id} />
            <ChapterActions id={chapter._id} />
          </div>
        </div>

        {/* Chapter Items */}
        <div className="text-center text-sm text-gray-500 font-medium py-3">
          {/* {JSON.stringify(chapterItems)} */}
          {chapterItems && chapterItems.length === 0 ? (
            <p>There are no items in this chapter yet</p>
          ) : (
            <DndContext
              sensors={sensors}
              collisionDetection={closestCenter}
              onDragStart={handleItemDragStart}
              onDragEnd={handleItemDragEnd}
            >
              <SortableContext
                items={
                  chapterItems ? chapterItems.map((item) => item.title) : []
                }
                strategy={verticalListSortingStrategy}
              >
                <div className="space-y-3">
                  {chapterItems &&
                    chapterItems.map((item) => (
                      <SortableItem
                        key={item.title}
                        // @ts-ignore
                        item={item}
                        chapterId={chapter._id}
                      />
                    ))}
                </div>
              </SortableContext>
              <DragOverlay>
                {/* @ts-ignore */}
                {activeItem && <SortableItem item={activeItem} />}
              </DragOverlay>
            </DndContext>
          )}
        </div>
      </div>
    </div>
  );
}

// Sortable Item Component
function SortableItem({
  item,
  chapterId,
}: {
  item: ICourseLesson;
  chapterId: number;
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: item._id });

  const style: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const Icon =
    item.type === "lesson"
      ? FiBookOpen
      : item.type === "quiz"
      ? MdChecklist
      : item.type === "assignment"
      ? MdOutlineFileCopy
      : null;

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={`w-full border rounded-sm p-3 bg-white cursor-grab ${
        "border-l-2 border-l-primary" //item.type === "Chapter" ? "" :
      }`}
    >
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          {Icon && <Icon className="text-primary" />}
          <span className="font-medium">{item.title}</span>
          <span className="ml-5 text-xs bg-gray-200 px-2 py-0.5 rounded-md">
            {item.is_published ? "Published" : "Draft"}
          </span>
        </div>
        <ItemActions
          chapterID={chapterId}
          itemID={item._id}
          itemType={item.type}
        />
      </div>
    </div>
  );
}
