"use client";

import Image from "next/image";
import React, { useState, useEffect } from "react";
import Link from "next/link";
import Checkbox from "@/components/ui/Checkbox";
import Pagination from "@/components/shared/Pagination";
import { formatDate } from "@/lib/helpers";

import {
  DndContext,
  closestCenter,
  DragEndEvent,
  useSensor,
  useSensors,
  PointerSensor,
  DragOverlay,
} from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { MdOutlineDragIndicator } from "react-icons/md";
import axiosInstance from "@/lib/axiosInstance";
import { useSession } from "next-auth/react";
import { toast } from "sonner";
import CoursesAction from "./CoursesAction";

interface CoursesListProps {
  data: CourseDetails[];
}

export default function CoursesList({ data }: CoursesListProps) {
  const { data: session } = useSession();

  const [courses, setCourses] = useState<CourseDetails[]>(data); // local copy
  const [selected, setSelected] = useState<number[]>([]);
  const [activeId, setActiveId] = useState<number | null>(null);

  const sensors = useSensors(useSensor(PointerSensor));

  useEffect(() => {
    setCourses(data); // sync props if they change
  }, [data]);

  const activeItem = courses.find((c) => c.id === activeId) || null;

  const toggleSelectAll = () => {
    if (selected.length === courses.length) {
      setSelected([]);
    } else {
      setSelected(courses.map((c) => c.id));
    }
  };

  const toggleSelectOne = (id: number) => {
    setSelected((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const isAllSelected =
    selected.length === courses.length && courses.length > 0;

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over) return;

    const activeId = active.id as number;
    const overId = over.id as number;

    const oldIndex = courses.findIndex((c) => c.id === activeId);
    const newIndex = courses.findIndex((c) => c.id === overId);

    if (oldIndex === -1 || newIndex === -1) return;

    // 1️⃣ Optimistic local reorder
    const newCourses = [...courses];
    const [moved] = newCourses.splice(oldIndex, 1);
    newCourses.splice(newIndex, 0, moved);
    setCourses(newCourses);

    // 2️⃣ Send to backend
    axiosInstance
      .put(
        `/private/course/reorder`,
        { activeID: activeId, overID: overId },
        { headers: { Authorization: `Bearer ${session?.accessToken}` } }
      )
      .catch(() => {
        toast.error("Reorder failed, reverting");
        // rollback
        setCourses(courses);
      });

    setActiveId(null);
  };

  return (
    <div className="border rounded-xl overflow-hidden">
      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragStart={(e) => setActiveId(e.active.id as number)}
        onDragEnd={handleDragEnd}
        onDragCancel={() => setActiveId(null)}
      >
        <table className="w-full text-sm">
          <thead className="bg-gray-100">
            <tr className="text-left">
              <th className="p-3">
                <div className="flex items-center gap-3 font-medium">
                  <Checkbox
                    checked={isAllSelected}
                    onChange={toggleSelectAll}
                  />
                  <span>Title</span>
                </div>
              </th>
              <th className="p-3 font-medium">Categories</th>
              <th className="p-3 font-medium">Author</th>
              <th className="p-3 font-medium">Price</th>
              <th className="p-3 font-medium">Date</th>
              <th className="p-3 font-medium">Actions</th>
            </tr>
          </thead>

          <SortableContext
            items={courses.map((c) => c.id)}
            strategy={verticalListSortingStrategy}
          >
            <tbody>
              {courses.map((course) => (
                <SortableRow
                  key={course.id}
                  course={course}
                  selected={selected}
                  toggleSelectOne={toggleSelectOne}
                />
              ))}

              {courses.length === 0 && (
                <tr>
                  <td colSpan={6} className="p-5 text-center text-gray-500">
                    No courses found.
                  </td>
                </tr>
              )}
            </tbody>
          </SortableContext>
        </table>

        <DragOverlay>
          {activeItem ? <SortableRow course={activeItem} isOverlay /> : null}
        </DragOverlay>
      </DndContext>

      {/* <Pagination totalPages={99} /> */}
    </div>
  );
}

// ---------------------------------------
// Sortable Row Component
// ---------------------------------------
function SortableRow({
  course,
  selected,
  toggleSelectOne,
  isOverlay = false,
}: {
  course: CourseDetails;
  selected?: number[];
  toggleSelectOne?: (id: number) => void;
  isOverlay?: boolean;
}) {
  const {
    setNodeRef,
    attributes,
    listeners,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: course.id });

  const style: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging && !isOverlay ? 0.4 : 1,
    background: isOverlay ? "#fff" : undefined,
    boxShadow: isOverlay ? "0 4px 12px rgba(0,0,0,0.15)" : undefined,
  };

  return (
    <tr
      ref={setNodeRef}
      style={style}
      {...(isOverlay ? {} : attributes)}
      className="border-t border-gray-300 hover:bg-gray-100"
    >
      <td className="p-3">
        <div className="flex items-center gap-3">
          <span {...listeners} className="text-gray-500 cursor-move">
            <MdOutlineDragIndicator />
          </span>

          {toggleSelectOne && (
            <Checkbox
              checked={selected?.includes(course.id) || false}
              onChange={() => toggleSelectOne(course.id)}
            />
          )}

          <div className="flex items-center gap-3">
            <Image
              src={course.featured_image || "/images/placeholder.svg"}
              alt="image"
              width={100}
              height={50}
              className="w-20 h-10 object-contain rounded-md"
            />
            <div>
              <p className="font-medium">{course.title}</p>
              <div className="flex items-center gap-2 mt-1">
                <p className="text-sm text-gray-600">
                  Topic:{" "}
                  <span className="text-black">
                    {course.course_chapters.length}
                  </span>
                </p>
                <p className="text-sm text-gray-600">
                  Lesson:{" "}
                  <span className="text-black">
                    {course.course_chapters.reduce(
                      (acc, chapter) => acc + chapter.course_lessons.length,
                      0
                    )}
                  </span>
                </p>
                <p className="text-sm text-gray-600">
                  Quiz: <span className="text-black">0</span>
                </p>
                <p className="text-sm text-gray-600">
                  Assignment: <span className="text-black">0</span>
                </p>
              </div>
            </div>
          </div>
        </div>
      </td>

      <td className="p-3">{course.general_settings.category.name}</td>

      <td className="p-3">
        <div className="flex items-center gap-2">
          <span className="bg-primary text-white w-8 h-8 rounded-full flex items-center justify-center font-medium">
            {course.author.name.slice(0, 1)}
          </span>
          {course.author.name}
        </div>
      </td>

      <td className="p-3">{course.sale_price}</td>

      <td className="p-3">
        <p>{formatDate(course.created_at)}</p>
      </td>

      <td className="p-3">
        <div className="flex items-center gap-2">
          {/* <Link
            href={`/courses/${course.id}`}
            className="text-primary border border-primary px-2 py-0.5 rounded-md hover:bg-primary hover:text-white"
          >
            View
          </Link> */}
          <CoursesAction id={course.id} />
        </div>
      </td>
    </tr>
  );
}
