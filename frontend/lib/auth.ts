import NextAuth from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";
import axiosInstance from "./axiosInstance";

export const { handlers, auth, signIn, signOut } = NextAuth({
  // Required behind Coolify/Traefik or any reverse proxy so callback URLs resolve correctly.
  trustHost: true,
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

      // If access token is missing or malformed, keep token and let the client handle re-auth.
      // Returning `null` here can lead to "Bearer undefined" headers and confusing UX.
      //@ts-ignore
      const accessToken: unknown = token.accessToken;
      if (typeof accessToken === "string") {
        try {
          const decodedAccessToken = JSON.parse(
            Buffer.from(accessToken.split(".")[1], "base64").toString()
          ) as { exp?: number };

          if (
            typeof decodedAccessToken?.exp === "number" &&
            decodedAccessToken.exp < Math.round(Date.now() / 1000)
          ) {
            //@ts-ignore
            token.error = "AccessTokenExpired";
          }
        } catch {
          //@ts-ignore
          token.error = "AccessTokenInvalid";
        }
      } else {
        //@ts-ignore
        token.error = "AccessTokenMissing";
      }

      return token;
    },
    session: async ({ session, token }) => {
      if (token) {
        //@ts-ignore
        session.accessToken = token.accessToken;
        //@ts-ignore
        session.user = token.user;
        //@ts-ignore
        session.error = token.error;
      }
      return session;
    },
  },
  pages: {
    signOut: "/",
  },
});
