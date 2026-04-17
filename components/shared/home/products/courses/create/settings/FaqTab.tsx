import InputField from "@/components/ui/InputField";
import React, { useState } from "react";
import { IoIosArrowDown, IoIosArrowUp } from "react-icons/io";
import { LuPlus, LuTrash2 } from "react-icons/lu";

type FaqItem = {
  id: number;
  question: string;
  answer: string;
};

export default function FaqTab() {
  const [faqs, setFaqs] = useState<FaqItem[]>([]);
  const [isOpen, setIsOpen] = useState(true);

  const handleCreateFaq = () => {
    const newFaq: FaqItem = {
      id: Date.now(),
      question: "",
      answer: "",
    };
    setFaqs((prev) => [...prev, newFaq]);
  };

  const handleChange = (id: number, field: keyof FaqItem, value: string) => {
    setFaqs((prev) =>
      prev.map((faq) => (faq.id === id ? { ...faq, [field]: value } : faq))
    );
  };

  const handleDeleteFaq = (itemId: number) => {
    setFaqs((prev) => prev.filter((faq) => faq.id !== itemId));
  };

  return (
    <div className="flex items-start gap-10">
      {/* Left side - FAQ editor */}
      <div className="w-full">
        <label className="block text-sm font-medium mb-2">FAQs</label>

        {faqs.map((item) => (
          <div
            key={item.id}
            className="relative w-full border rounded-xl p-5 mb-4"
          >
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">Question</label>
              <InputField
                placeholder="Write your question"
                value={item.question}
                onChange={(e) =>
                  handleChange(item.id, "question", e.target.value)
                }
                className="w-full"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">Answer</label>
              <textarea
                rows={4}
                className="bg-white text-sm border w-full rounded-md px-3 py-2 outline-none focus:border-primary"
                placeholder="Write your answer"
                value={item.answer}
                onChange={(e) =>
                  handleChange(item.id, "answer", e.target.value)
                }
              />
            </div>
            <button
              onClick={() => handleDeleteFaq(item.id)}
              className="absolute bg-red-200 p-1 rounded-full -top-1.5 -right-1.5 text-red-500"
            >
              <LuTrash2 size={14} />
            </button>
          </div>
        ))}

        <button
          onClick={handleCreateFaq}
          className="px-4 py-2 rounded-md text-primary border border-primary text-sm font-medium flex items-center gap-2"
        >
          <LuPlus /> Add New
        </button>
      </div>

      {/* Right side - Info box */}
      <div className="w-80 min-w-80">
        <div className="border p-5 rounded-md">
          <button
            onClick={() => setIsOpen(!isOpen)}
            className={`flex items-center justify-between w-full ${
              isOpen ? "text-primary" : ""
            }`}
          >
            <p className="font-semibold text-left">
              Why should I mention the language?
            </p>
            {isOpen ? <IoIosArrowUp /> : <IoIosArrowDown />}
          </button>

          {isOpen && (
            <p className="text-sm mt-5 text-gray-600">
              It is important to mention the language of the course content so
              that students who understand the language can enroll in the
              course.
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
