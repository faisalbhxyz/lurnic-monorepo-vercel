/**
 * Regression checks for course/lesson schedule_time helpers (run: npx tsx scripts/verify-schedule-time.ts).
 */
import { hhmmToScheduleTime, scheduleTimeToHHMM } from "../lib/helpers";

function assert(cond: boolean, msg: string): void {
  if (!cond) throw new Error(`FAIL: ${msg}`);
}

assert(scheduleTimeToHHMM(null) === "", "null -> empty");
assert(scheduleTimeToHHMM("") === "", "empty -> empty");
assert(scheduleTimeToHHMM("10:05 AM") === "10:05", "12h AM");
assert(scheduleTimeToHHMM("10:05 PM") === "22:05", "12h PM");
assert(scheduleTimeToHHMM("12:00 AM") === "00:00", "midnight");
assert(scheduleTimeToHHMM("12:30 PM") === "12:30", "noon+30");
assert(scheduleTimeToHHMM("10:05:30") === "10:05", "DB HH:MM:SS");

assert(hhmmToScheduleTime("10:05") === "10:05 AM", "input -> AM");
assert(hhmmToScheduleTime("22:05") === "10:05 PM", "input -> PM");
assert(hhmmToScheduleTime("00:30") === "12:30 AM", "00:30 -> 12:30 AM");
assert(hhmmToScheduleTime("12:00") === "12:00 PM", "12:00 -> noon");

console.log("verify-schedule-time: OK");
