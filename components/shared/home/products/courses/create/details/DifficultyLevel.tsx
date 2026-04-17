import {
  Listbox,
  ListboxButton,
  ListboxOption,
  ListboxOptions,
} from "@headlessui/react";
import { IoIosArrowDown } from "react-icons/io";
import clsx from "clsx";
import { useState } from "react";

const people = [
  { id: 1, name: "All Levels" },
  { id: 2, name: "Beginner" },
  { id: 3, name: "Intermediate" },
  { id: 3, name: "Expert" },
];

export default function DifficultyLevel() {
  const [selected, setSelected] = useState(people[0]);
  return (
    <Listbox value={selected} onChange={setSelected}>
      <ListboxButton
        className={clsx(
          "relative flex items-center justify-between border w-full rounded-lg bg-white py-1.5 px-3 text-left text-sm",
          "focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25"
        )}
      >
        {selected.name}
        <IoIosArrowDown size={16} className="text-gray-500" />
      </ListboxButton>
      <ListboxOptions
        anchor="bottom"
        transition
        className={clsx(
          "w-[var(--button-width)] rounded-xl border bg-white p-1 [--anchor-gap:var(--spacing-1)] focus:outline-none",
          "transition duration-100 ease-in data-[leave]:data-[closed]:opacity-0"
        )}
      >
        {people.map((person) => (
          <ListboxOption
            key={person.name}
            value={person}
            className="flex cursor-default items-center gap-2 rounded-lg py-1.5 px-3 select-none data-[focus]:bg-gray-100 data-[selected]:bg-gray-100"
          >
            <div className="text-sm/6">{person.name}</div>
          </ListboxOption>
        ))}
      </ListboxOptions>
    </Listbox>
  );
}
