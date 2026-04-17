import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { HiDotsVertical } from "react-icons/hi";
import { BiEditAlt } from "react-icons/bi";
import { HiOutlineTrash } from "react-icons/hi";

export default function CountryActions() {
  return (
    <Menu>
      <MenuButton className="inline-flex items-center gap-2 rounded-md p-1.5 text-sm font-semibold shadow-inner shadow-white/10 focus:outline-none data-[hover]:bg-gray-100 data-[open]:bg-gray-100 data-[focus]:outline-1 data-[focus]:outline-white">
        <HiDotsVertical size={18} className="text-gray-600" />
      </MenuButton>
      <MenuItems
        transition
        anchor="bottom end"
        className="w-40 origin-top-right rounded-lg border bg-white p-1 text-sm transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
      >
        <MenuItem>
          <button className="group flex w-full items-center gap-2 rounded-lg py-1.5 px-3 data-[focus]:bg-gray-100">
            <BiEditAlt size={18} />
            Edit
          </button>
        </MenuItem>
        <MenuItem>
          <button className="group flex w-full items-center text-red-500 gap-2 rounded-lg py-1.5 px-3 data-[focus]:bg-gray-100">
            <HiOutlineTrash size={18} />
            Delete
          </button>
        </MenuItem>
      </MenuItems>
    </Menu>
  );
}
