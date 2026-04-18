import { hhmmToScheduleTime, scheduleTimeToHHMM } from "@/lib/helpers";
import { LuClock } from "react-icons/lu";
import { Control, Controller, FieldValues, Path } from "react-hook-form";

type ScheduleTimeFieldProps<T extends FieldValues> = {
  control: Control<T>;
  name: Path<T>;
};

/**
 * Native time input (minute precision, keyboard-friendly) — values are stored as
 * "03:04 PM" strings for the Go API (`time.Parse("03:04 PM", …)`).
 */
export default function ScheduleTimeField<T extends FieldValues>({
  control,
  name,
}: ScheduleTimeFieldProps<T>) {
  return (
    <Controller
      control={control}
      name={name}
      render={({ field: { value, onChange } }) => (
        <div className="w-full flex items-center justify-center gap-2 text-sm py-2 px-1">
          <LuClock size={18} className="text-gray-500 shrink-0" aria-hidden />
          <input
            type="time"
            step={60}
            title="Schedule time"
            className="w-full min-w-0 flex-1 bg-transparent border-0 p-0 text-sm focus:outline-none focus:ring-0 [color-scheme:light]"
            value={scheduleTimeToHHMM(value)}
            onChange={(e) => {
              const v = e.target.value;
              if (!v) {
                onChange(null);
                return;
              }
              onChange(hhmmToScheduleTime(v));
            }}
          />
        </div>
      )}
    />
  );
}
