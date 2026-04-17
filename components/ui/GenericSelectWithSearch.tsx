import React, { useState, useRef, useEffect } from "react";
import { IoIosArrowDown } from "react-icons/io";

type GenericSelectWithSearchProps<T> = {
  items: T[];
  getLabel: (item: T) => string;
  onSelect: (item: T) => void;
  selectedItem?: T | null;
  placeholder?: string;
  className?: string;
};

export function GenericSelectWithSearch<T>({
  items,
  getLabel,
  onSelect,
  selectedItem = null,
  placeholder = "Choose one",
  className = "",
}: GenericSelectWithSearchProps<T>) {
  const [isOpen, setIsOpen] = useState(false);
  const [search, setSearch] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);

  const filteredItems = items.filter((item) =>
    getLabel(item).toLowerCase().includes(search.toLowerCase())
  );

  function toggleOpen() {
    setIsOpen((prev) => !prev);
  }

  function handleSelect(item: T) {
    onSelect(item);
    setIsOpen(false);
    setSearch("");
  }

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    }

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    } else {
      document.removeEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isOpen]);

  return (
    <div ref={containerRef} className={`relative w-60 ${className}`}>
      <button
        type="button"
        onClick={toggleOpen}
        className="w-full text-sm flex border items-center justify-between px-3 py-1.5 rounded-md"
      >
        <span>{selectedItem ? getLabel(selectedItem) : placeholder}</span>
        <IoIosArrowDown />
      </button>

      {isOpen && (
        <div className="absolute top-full border mt-1 left-0 right-0 bg-white rounded-lg p-2 z-10">
          <input
            type="text"
            placeholder="Search..."
            className="border w-full text-sm px-3 py-1.5 rounded-md mb-2"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            autoFocus
          />
          <ul className="max-h-48 overflow-y-auto">
            {filteredItems.length === 0 && (
              <li className="p-2 text-sm text-gray-500">No results found</li>
            )}
            {filteredItems.map((item, i) => (
              <li
                key={i}
                onClick={() => handleSelect(item)}
                className="cursor-pointer px-2 py-1 hover:bg-gray-100 rounded text-sm"
              >
                {getLabel(item)}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}
