import React, { useState } from "react";
import Checkbox from "@/components/ui/Checkbox";
import { Menu, MenuButton, MenuItems } from "@headlessui/react";
import { IoIosArrowDown } from "react-icons/io";

type Category = {
  id: number;
  name: string;
  subCategory: Category[];
};

export default function Categories() {
  const [categories, setCategories] = useState<Category[]>([
    {
      id: 1,
      name: "category 1",
      subCategory: [{ id: 1, name: "sub 11", subCategory: [] }],
    },
    {
      id: 2,
      name: "category 2",
      subCategory: [{ id: 1, name: "sub 21", subCategory: [] }],
    },
  ]);

  const [newCategoryName, setNewCategoryName] = useState("");
  const [selected, setSelected] = useState<Category | null>(null);

  const handleAddCategory = () => {
    if (!newCategoryName.trim()) return;

    const newCategory: Category = {
      id: Date.now(),
      name: newCategoryName,
      subCategory: [],
    };

    if (selected) {
      setCategories((prev) =>
        prev.map((cat) =>
          cat.id === selected.id
            ? { ...cat, subCategory: [...cat.subCategory, newCategory] }
            : cat
        )
      );
    } else {
      setCategories((prev) => [...prev, newCategory]);
    }

    setNewCategoryName("");
    setSelected(null);
  };

  const renderCategory = (category: Category, level = 0) => (
    <li key={category.id} className="flex items-center gap-2">
      <div className={`ml-${level * 4}`}>
        <Checkbox />
        <span className="text-sm ml-2">{category.name}</span>
        {category.subCategory.length > 0 && (
          <ul className="mt-1 space-y-1">
            {category.subCategory.map((sub) => renderCategory(sub, level + 1))}
          </ul>
        )}
      </div>
    </li>
  );

  const renderSelectableCategory = (category: Category, level = 0) => (
    <li
      key={category.id}
      className={`mt-1 ml-${level * 4} cursor-pointer px-2 py-1 rounded`}
      onClick={() => setSelected(category)}
    >
      <span className="text-sm">{category.name}</span>
      {category.subCategory.length > 0 && (
        <ul className="ml-4">
          {category.subCategory.map((sub) =>
            renderSelectableCategory(sub, level + 1)
          )}
        </ul>
      )}
    </li>
  );

  return (
    <div className="border rounded-md bg-white p-3">
      <ul className="mb-3 space-y-2">
        {categories.map((category) => renderCategory(category))}
      </ul>

      <Menu>
        <MenuButton className="text-sm text-primary font-medium">
          + Add
        </MenuButton>
        <MenuItems
          transition
          anchor="bottom start"
          className="w-80 p-3 origin-top-left rounded-md border bg-white text-sm transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
        >
          <div className="mb-5">
            <input
              type="text"
              placeholder="Category name"
              value={newCategoryName}
              onChange={(e) => setNewCategoryName(e.target.value)}
              className="border outline-none w-full px-3 py-1.5 rounded-md"
            />
          </div>

          <div>
            <button className="w-full border px-3 py-1.5 rounded-md flex items-center justify-between">
              {selected?.name ?? "Select parent"}
              <IoIosArrowDown />
            </button>
            <div className="p-2">
              <ul>{categories.map((cat) => renderSelectableCategory(cat))}</ul>
            </div>
          </div>

          <div className="flex items-center justify-end gap-2 mt-3">
            <button
              onClick={() => {
                setNewCategoryName("");
                setSelected(null);
              }}
              className="text-gray-500 text-sm"
            >
              Cancel
            </button>
            <button
              onClick={handleAddCategory}
              className="bg-primary/10 text-primary px-3 py-1 rounded-md"
            >
              Ok
            </button>
          </div>
        </MenuItems>
      </Menu>
    </div>
  );
}
