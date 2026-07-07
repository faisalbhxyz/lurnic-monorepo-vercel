import {
  buildAttachmentPreviews,
  normalizeAssignmentAttachments,
  revokeBlobUrls,
} from "../lib/assignmentAttachments";

function assert(condition: boolean, message: string) {
  if (!condition) {
    throw new Error(message);
  }
}

function testNormalizeStoredShape() {
  const raw = [
    {
      url: "https://cdn.example.com/brief.pdf",
      file_name: "brief.pdf",
      mime_type: "application/pdf",
      size: 2048,
    },
  ];
  const normalized = normalizeAssignmentAttachments(raw);
  assert(normalized?.length === 1, "expected one stored attachment");
  assert(normalized?.[0].url === raw[0].url, "url mismatch");
  assert(normalized?.[0].file_name === raw[0].file_name, "file_name mismatch");
}

function testNormalizeLooseDbShape() {
  const raw = [{ file_path: "https://cdn.example.com/a.png", name: "a.png", type: "image/png" }];
  const normalized = normalizeAssignmentAttachments(raw);
  assert(normalized?.[0].url === "https://cdn.example.com/a.png", "loose url mismatch");
  assert(normalized?.[0].file_name === "a.png", "loose name mismatch");
}

function testNormalizeJsonString() {
  const raw = JSON.stringify([
    { url: "https://cdn.example.com/x.pdf", file_name: "x.pdf" },
  ]);
  const normalized = normalizeAssignmentAttachments(raw);
  assert(normalized?.length === 1, "json string should parse to one attachment");
}

function testBuildPreviewsForFileAndStored() {
  const file = new File(["hello"], "note.txt", { type: "text/plain" });
  const stored = {
    url: "https://cdn.example.com/img.png",
    file_name: "img.png",
    mime_type: "image/png",
  };
  const { previews, blobUrls } = buildAttachmentPreviews([stored, file]);
  assert(previews.length === 2, "expected two previews");
  assert(previews[0].name === "img.png", "stored preview name mismatch");
  assert(previews[1].name === "note.txt", "file preview name mismatch");
  revokeBlobUrls(blobUrls);
}

function main() {
  testNormalizeStoredShape();
  testNormalizeLooseDbShape();
  testNormalizeJsonString();
  testBuildPreviewsForFileAndStored();
  console.log("verify-assignment-attachments: all checks passed");
}

main();
