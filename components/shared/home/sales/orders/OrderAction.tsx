import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { HiDotsVertical } from "react-icons/hi";
import { BiEditAlt } from "react-icons/bi";
import { useRouter } from "next/navigation";
import { useSession } from "next-auth/react";
import axiosInstance from "@/lib/axiosInstance";
import { toast } from "sonner";

export default function OrderAction({
  isPaid,
  id,
}: {
  isPaid: boolean;
  id: number;
}) {
  const router = useRouter();
  const { data: session } = useSession();

  const handleDelete = () => {
    const isConfirm = confirm("Are you sure you want to delete this banner?");
    if (!isConfirm) return;
    axiosInstance
      .delete(`/private/order/delete/${id}`, {
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

  const handleMarkAsPaid = (id: number) => {
    const isConfirm = confirm(
      "Are you sure you want to mark this order as paid?"
    );
    if (!isConfirm) return;
    axiosInstance
      .put(
        `/private/order/mark-as-paid/${id}`,
        {},
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${session?.accessToken}`,
          },
        }
      )
      .then((res) => {
        toast.success(res.data.message);
        router.refresh();
      })
      .catch((error) => {
        console.log("[ERROR]", error);

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
        {!isPaid && (
          <MenuItem>
            <button
              className="group flex w-full items-center gap-2 rounded-lg py-1.5 px-3 data-[focus]:bg-blue-500 data-[focus]:text-white disabled:cursor-not-allowed disabled:opacity-50"
              onClick={() => handleMarkAsPaid(id)}
              disabled={isPaid}
            >
              <BiEditAlt size={18} />
              Mark As Paid
            </button>
          </MenuItem>
        )}
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
