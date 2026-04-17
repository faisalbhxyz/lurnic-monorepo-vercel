import NextAuth from "next-auth";
import { JWT } from "next-auth/jwt";

declare module "next-auth" {
  interface Session {
    accessToken: string;
    refreshToken: string;
    error?: string;
    user: {
      user_id: string;
      name: string;
      email: string;
      phone: string;
      role: string;
      permissions: string[];
      // subscriptionToken: string;
      // subscriptionIAT: number;
      // subscriptionExpiry: number;
    };
  }
}
