import { auth } from "./lib/auth";

export default auth((req) => {
  // const isHomePage = req.nextUrl.pathname === "/";
  const isLoginPage = req.nextUrl.pathname === "/login";
  const isSignupPage = req.nextUrl.pathname === "/signup";

  // If the user is not authenticated and is not on the login or signup page
  if (!req.auth && !isLoginPage && !isSignupPage) {
    const newUrl = new URL("/login", req.nextUrl.origin);
    return Response.redirect(newUrl);
  }
  if (req.auth && (isLoginPage || isSignupPage)) {
    const newUrl = new URL("/", req.nextUrl.origin);
    return Response.redirect(newUrl);
  }
});

export const config = {
  // Exclude /v1 so rewrites to the Go API (see root vercel.json) are not intercepted by NextAuth.
  matcher: ["/((?!api|v1|_next/static|_next/image|favicon.ico).*)"],
};
