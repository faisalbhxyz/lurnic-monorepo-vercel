export const formatDate = (date: string): string => {
  const parsedDate = new Date(date);
  const day = String(parsedDate.getDate()).padStart(2, "0");
  const month = parsedDate.toLocaleString("en-US", { month: "long" }); // Full month name
  const year = parsedDate.getFullYear();

  return `${day} ${month}, ${year}`;
};

export const dbTimeToPickerFormat = (timeStr: string): string => {
  const [hourStr, minStr] = timeStr.split(":");
  let hour = parseInt(hourStr, 10);
  const period = hour < 12 ? "AM" : "PM";
  hour = hour % 12 === 0 ? 12 : hour % 12;
  return `${hour.toString().padStart(2, "0")}:${minStr} ${period}`;
};
