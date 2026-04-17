import {
  Listbox,
  ListboxButton,
  ListboxOption,
  ListboxOptions,
} from "@headlessui/react";
import { IoIosArrowDown } from "react-icons/io";
import clsx from "clsx";
import { IoEyeSharp } from "react-icons/io5";
import { Controller, useFormContext } from "react-hook-form";
import { TCourseSchema } from "@/schema/course.schema";

const people = [
  { id: 1, name: "Public", value: "public" },
  // { id: 2, name: "Password Protected", value: "protected" },
  { id: 3, name: "Private", value: "private" },
];

export default function Visibility() {
  const { control } = useFormContext<TCourseSchema>();

  return (
    <Controller
      control={control}
      name="visibility"
      render={({ field: { onChange, value } }) => {
        const selected = people.find((person) => person.value === value);

        return (
          <Listbox value={value} onChange={onChange}>
            <ListboxButton
              className={clsx(
                "relative flex items-center justify-between border w-full rounded-lg bg-white py-1.5 px-3 text-left text-sm",
                "focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25"
              )}
            >
              <div className="flex items-center gap-1">
                <IoEyeSharp size={20} className="text-gray-500" />
                {selected?.name ?? "Select visibility"}
              </div>
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
                  key={person.id}
                  value={person.value} // ✅ use the string value here
                  className="flex cursor-default items-center gap-2 rounded-lg py-1.5 px-3 select-none data-[focus]:bg-gray-100 data-[selected]:bg-gray-100"
                >
                  <div className="text-sm/6">{person.name}</div>
                </ListboxOption>
              ))}
            </ListboxOptions>
          </Listbox>
        );
      }}
    />
  );
}
