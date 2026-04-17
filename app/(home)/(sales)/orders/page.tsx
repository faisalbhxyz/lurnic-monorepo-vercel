import Orders from "@/components/shared/home/sales/orders/Orders";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import React from "react";

const getAllOrders = async (session: Session) => {
  try {
    const res = await axiosInstance.get("/private/order", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });

    return res.data.data;
  } catch (error) {
    return [];
  }
};

export default async function page() {
  const session = await auth();
  if (!session) return null;

  const orders = await getAllOrders(session);

  // console.log("orders", orders);
  

  return <Orders orders={orders}/>;
}
