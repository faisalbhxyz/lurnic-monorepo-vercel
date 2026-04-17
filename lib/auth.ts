import NextAuth from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";
import axiosInstance from "./axiosInstance";

export const { handlers, auth, signIn, signOut } = NextAuth({
  session: {
    strategy: "jwt",
    maxAge: 2 * 60 * 60,
  },
  providers: [
    CredentialsProvider({
      async authorize(credentials) {
        if (credentials === null) return null;
        try {
          const res = await axiosInstance.post("/user/login", {
            email: credentials.email,
            password: credentials.password,
          });

          const authInfo = res.data;

          if (res.status === 200 && authInfo) {
            return {
              token: authInfo.token,
              user: authInfo.user,
            } as any;
          } else {
            throw new Error("Invalid credentials");
          }
        } catch (error: any) {
          throw error;
        }
      },
    }),
  ],
  jwt: {
    maxAge: 2 * 60 * 60,
  },
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        //@ts-ignore
        token.accessToken = user.token;
        //@ts-ignore
        token.user = user.user;
      }

      const decodedAccessToken = JSON.parse(
        //@ts-ignore
        Buffer.from(token.accessToken.split(".")[1], "base64").toString()
      );

      if (decodedAccessToken.exp < Math.round(Date.now() / 1000)) {
        return null;
      }

      return token;
    },
    session: async ({ session, token }) => {
      if (token) {
        //@ts-ignore
        session.accessToken = token.accessToken;
        //@ts-ignore
        session.user = token.user;
      }
      return session;
    },
  },
  pages: {
    signOut: "/",
  },
});
