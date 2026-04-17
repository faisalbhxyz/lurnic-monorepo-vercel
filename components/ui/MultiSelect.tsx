import React, { useEffect, useState } from "react";
import { IoIosArrowDown } from "react-icons/io";
import { IoCheckmark } from "react-icons/io5";

export interface Option {
  id: number | string;
  name: string;
}

interface MultiSelectProps {
  options: Option[];
  selected: Option[];
  onChange: (selected: Option[]) => void;
  placeholder?: string;
}

const MultiSelect: React.FC<MultiSelectProps> = ({
  options,
  selected,
  onChange,
  placeholder = "Search...",
}) => {
  const [searchValue, setSearchValue] = useState<string>("");
  const [isOpenDropdown, setIsOpenDropdown] = useState<boolean>(false);

  const filteredItems = options.filter((item) =>
    item.name.toLowerCase().includes(searchValue.toLowerCase())
  );

  const isSelected = (item: Option): boolean =>
    selected.some((selectedItem) => selectedItem.id === item.id);

  const toggleSelectItem = (item: Option): void => {
    if (isSelected(item)) {
      onChange(selected.filter((selectedItem) => selectedItem.id !== item.id));
    } else {
      onChange([...selected, item]);
    }
  };

  const removeItem = (option: Option): void => {
    onChange(selected.filter((selectedItem) => selectedItem.id !== option.id));
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as HTMLElement;
      if (!target.closest(".custom-select")) {
        setIsOpenDropdown(false);
      }
    };

    document.addEventListener("click", handleClickOutside);
    return () => {
      document.removeEventListener("click", handleClickOutside);
    };
  }, []);

  return (
    <div className="relative custom-select w-full">
      <input
        type="text"
        placeholder={placeholder}
        value={searchValue}
        onChange={(e) => setSearchValue(e.target.value)}
        onFocus={() => setIsOpenDropdown(true)}
        className="w-full text-sm bg-white border rounded-md px-3 py-2 focus:outline-none"
      />

      <IoIosArrowDown
        className={`${
          isOpenDropdown ? "rotate-180" : "rotate-0"
        } transition-all duration-300 text-[1rem] absolute top-[10px] right-3 text-gray-500`}
      />

      {isOpenDropdown && (
        <div className="absolute border bg-white left-0 w-full mt-1 rounded-md shadow-lg z-20 max-h-60 overflow-auto">
          {filteredItems.map((item) => (
            <div
              key={item.id}
              onClick={() => toggleSelectItem(item)}
              className="cursor-pointer text-gray-700 hover:bg-gray-100 px-3 py-2 flex items-center justify-between text-sm"
            >
              {item.name}
              <IoCheckmark
                className={`${
                  isSelected(item)
                    ? "scale-100 opacity-100"
                    : "scale-50 opacity-0"
                } mr-2 transition-all text-[1.3rem] duration-300`}
              />
            </div>
          ))}
          {filteredItems.length === 0 && (
            <p className="text-center text-gray-500 py-4 text-sm">
              No search results!
            </p>
          )}
        </div>
      )}

      {selected.length > 0 && (
        <div className="flex items-center gap-2 flex-wrap mt-2">
          {selected.map((item) => (
            <div
              key={item.id}
              className="bg-primary/20 text-primary px-3 py-1 rounded-full flex items-center text-sm"
            >
              {item.name}
              <button onClick={() => removeItem(item)} className="ml-2">
                ×
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default MultiSelect;
