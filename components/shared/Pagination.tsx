import React, { useState } from "react";
import { MdKeyboardArrowLeft, MdKeyboardArrowRight } from "react-icons/md";
import { HiOutlineDotsHorizontal } from "react-icons/hi";

interface PaginationProps {
  totalPages: number;
  perPageOptions?: number[];
  defaultPerPage?: number;
}

export default function Pagination({
  totalPages,
  perPageOptions = [10, 15, 20],
  defaultPerPage = 10,
}: PaginationProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const [perPage, setPerPage] = useState(defaultPerPage);

  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page);
    }
  };

  const handlePerPageChange = (value: number) => {
    setPerPage(value);
    setCurrentPage(1); // Reset to first page on perPage change
  };

  const renderPageButtons = () => {
    const pages: (number | string)[] = [];

    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);
      if (currentPage > 4) pages.push("dots-left");
      for (
        let i = Math.max(2, currentPage - 1);
        i <= Math.min(totalPages - 1, currentPage + 1);
        i++
      ) {
        pages.push(i);
      }
      if (currentPage < totalPages - 3) pages.push("dots-right");
      pages.push(totalPages);
    }

    return pages.map((page, index) =>
      typeof page === "number" ? (
        <button
          key={index}
          onClick={() => handlePageChange(page)}
          className={`px-3 py-1 border rounded ${
            currentPage === page
              ? "bg-gray-300 font-semibold"
              : "hover:bg-gray-100"
          }`}
        >
          {page}
        </button>
      ) : (
        <span key={index} className="px-2 text-gray-500">
          <HiOutlineDotsHorizontal />
        </span>
      )
    );
  };

  return (
    <div className="flex flex-col md:flex-row items-center justify-between gap-4 py-4 px-5 border-t border-gray-300 text-sm">
      {/* Per Page Selector */}
      <div className="flex items-center gap-2">
        <span>Show</span>
        <select
          value={perPage}
          onChange={(e) => handlePerPageChange(parseInt(e.target.value))}
          className="border rounded px-2 py-1"
        >
          {perPageOptions.map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
        <span>per page</span>
      </div>

      {/* Pagination Buttons */}
      <div className="flex items-center gap-1">
        <button
          onClick={() => handlePageChange(currentPage - 1)}
          disabled={currentPage === 1}
          className="p-1 border rounded hover:bg-gray-100 disabled:opacity-50"
        >
          <MdKeyboardArrowLeft size={20} />
        </button>
        {renderPageButtons()}
        <button
          onClick={() => handlePageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          className="p-1 border rounded hover:bg-gray-100 disabled:opacity-50"
        >
          <MdKeyboardArrowRight size={20} />
        </button>
      </div>
    </div>
  );
}
