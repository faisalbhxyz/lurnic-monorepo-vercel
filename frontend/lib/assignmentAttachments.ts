export type StoredAttachment = {
  url: string;
  file_name: string;
  mime_type?: string;
  size?: number;
};

export type AttachmentPreview = {
  type: "image" | "file";
  src: string;
  name: string;
};

export function isStoredAttachment(value: unknown): value is StoredAttachment {
  return (
    typeof value === "object" &&
    value !== null &&
    "url" in value &&
    typeof (value as StoredAttachment).url === "string" &&
    "file_name" in value &&
    typeof (value as StoredAttachment).file_name === "string"
  );
}

function normalizeAttachmentItem(item: unknown): StoredAttachment | null {
  if (item instanceof File) {
    return null;
  }
  if (isStoredAttachment(item)) {
    return item;
  }
  if (typeof item === "object" && item !== null) {
    const obj = item as Record<string, unknown>;
    const url = obj.url ?? obj.file_path ?? obj.URL;
    const file_name = obj.file_name ?? obj.name ?? obj.FileName;
    if (typeof url === "string" && typeof file_name === "string") {
      return {
        url,
        file_name,
        mime_type:
          typeof obj.mime_type === "string"
            ? obj.mime_type
            : typeof obj.type === "string"
              ? obj.type
              : typeof obj.MimeType === "string"
                ? obj.MimeType
                : undefined,
        size:
          typeof obj.size === "number"
            ? obj.size
            : typeof obj.Size === "number"
              ? obj.Size
              : undefined,
      };
    }
  }
  return null;
}

export function normalizeAssignmentAttachments(
  raw: unknown
): StoredAttachment[] | null {
  if (raw == null) return null;

  let items: unknown[] = [];
  if (typeof raw === "string") {
    try {
      const parsed = JSON.parse(raw);
      items = Array.isArray(parsed) ? parsed : [];
    } catch {
      return null;
    }
  } else if (Array.isArray(raw)) {
    items = raw;
  } else {
    return null;
  }

  const normalized = items
    .map(normalizeAttachmentItem)
    .filter((item): item is StoredAttachment => item !== null);

  return normalized.length > 0 ? normalized : null;
}

export function buildAttachmentPreviews(
  attachments: unknown[] | null | undefined
): { previews: AttachmentPreview[]; blobUrls: string[] } {
  const blobUrls: string[] = [];
  const previews: AttachmentPreview[] = [];

  for (const item of attachments ?? []) {
    if (item instanceof File) {
      if (item.type.startsWith("image/")) {
        const src = URL.createObjectURL(item);
        blobUrls.push(src);
        previews.push({ type: "image", src, name: item.name });
      } else {
        previews.push({ type: "file", src: "", name: item.name });
      }
      continue;
    }

    const stored = normalizeAttachmentItem(item);
    if (stored) {
      previews.push({
        type: stored.mime_type?.startsWith("image/") ? "image" : "file",
        src: stored.url,
        name: stored.file_name,
      });
    }
  }

  return { previews, blobUrls };
}

export function revokeBlobUrls(urls: string[]) {
  urls.forEach((url) => {
    if (url.startsWith("blob:")) {
      URL.revokeObjectURL(url);
    }
  });
}
