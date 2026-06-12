import Link from "next/link";

export default function NotFound() {
  return (
    <div className="flex min-h-[50vh] flex-col items-center justify-center gap-4">
      <h1 className="text-2xl font-semibold">Page not found</h1>
      <p className="text-gray-500">
        The page you are looking for does not exist.
      </p>
      <Link href="/" className="text-[#00828a] hover:underline">
        Go to dashboard
      </Link>
    </div>
  );
}
