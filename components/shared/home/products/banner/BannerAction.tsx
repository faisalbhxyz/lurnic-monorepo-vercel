import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { HiDotsVertical } from "react-icons/hi";
import { BiEditAlt } from "react-icons/bi";
import { useRouter } from "next/navigation";
import { useSession } from "next-auth/react";
import axiosInstance from "@/lib/axiosInstance";
import { toast } from "sonner";

export default function BannerAction({ id }: { id: number }) {
  const router = useRouter();
  const { data: session } = useSession();

  const handleDelete = () => {
    const isConfirm = confirm("Are you sure you want to delete this banner?");
    if (!isConfirm) return;
    axiosInstance
      .delete(`/private/banner/delete/${id}`, {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${session?.accessToken}`,
        },
      })
      .then((res) => {
        toast.success(res.data.message);
        router.refresh();
      })
      .catch((error) => {
        toast.error(error.response.data.error || "Something went wrong.");
      });
  };

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
          <button
            className="group flex w-full items-center gap-2 rounded-lg py-1.5 px-3 data-[focus]:bg-blue-500 data-[focus]:text-white"
            onClick={() => router.push(`/banner/${id}/edit`)}
          >
            <BiEditAlt size={18} />
            Edit
          </button>
        </MenuItem>
        <MenuItem>
          <button
            className="group flex w-full items-center gap-2 rounded-lg py-1.5 px-3 data-[focus]:bg-red-500 data-[focus]:text-white"
            onClick={handleDelete}
          >
            <BiEditAlt size={18} />
            Delete
          </button>
        </MenuItem>
      </MenuItems>
    </Menu>
  );
}
