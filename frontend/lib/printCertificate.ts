/**
 * Opens a print dialog with only the certificate preview so users can save as PDF.
 * Uses a hidden iframe (not window.open) so pop-up blockers do not interfere.
 * Clones the on-screen certificate and forces full A4 landscape dimensions so
 * the printed PDF matches the admin preview layout (no preview scale transform).
 */
export function printCertificatePreview(sourceEl: HTMLElement): void {
  const iframe = document.createElement("iframe");
  iframe.setAttribute("aria-hidden", "true");
  iframe.style.cssText =
    "position:fixed;width:0;height:0;border:0;visibility:hidden;";
  document.body.appendChild(iframe);

  const cleanup = () => {
    iframe.remove();
  };

  const win = iframe.contentWindow;
  if (!win) {
    cleanup();
    return;
  }

  const clone = sourceEl.cloneNode(true) as HTMLElement;
  prepareCertificateCloneForPrint(clone);
  stripNonPrintableElements(clone);

  const headStyles = Array.from(
    document.querySelectorAll<HTMLLinkElement | HTMLStyleElement>(
      'link[rel="stylesheet"], style'
    )
  )
    .map((el) => el.outerHTML)
    .join("\n");

  win.document.open();
  win.document.write(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>Certificate</title>
  ${headStyles}
  <style>
    @page { size: A4 landscape; margin: 0; }
    html, body {
      margin: 0;
      padding: 0;
      background: #ffffff;
    }
    body {
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      overflow: hidden;
    }
    .certificate-print-root {
      width: 297mm;
      height: 210mm;
      position: relative;
      overflow: hidden;
      background: #ffffff;
    }
    @media print {
      html, body {
        width: 297mm;
        height: 210mm;
        overflow: hidden;
      }
      .certificate-print-root {
        width: 297mm;
        height: 210mm;
      }
    }
  </style>
</head>
<body>
  <div class="certificate-print-root">${clone.innerHTML}</div>
</body>
</html>`);
  win.document.close();

  const printWhenReady = () => {
    win.focus();
    win.print();
    win.addEventListener("afterprint", cleanup, { once: true });
    // Fallback when afterprint is not fired (some browsers / cancel).
    setTimeout(cleanup, 60_000);
  };

  const images = Array.from(win.document.images);
  if (images.length === 0) {
    printWhenReady();
    return;
  }

  let loaded = 0;
  const onImageSettled = () => {
    loaded += 1;
    if (loaded >= images.length) {
      printWhenReady();
    }
  };

  images.forEach((img) => {
    if (img.complete) {
      onImageSettled();
    } else {
      img.addEventListener("load", onImageSettled, { once: true });
      img.addEventListener("error", onImageSettled, { once: true });
    }
  });
}

function stripNonPrintableElements(root: HTMLElement) {
  root.querySelectorAll("[data-print-exclude]").forEach((el) => el.remove());
}

function prepareCertificateCloneForPrint(root: HTMLElement) {
  root.querySelectorAll<HTMLElement>("[style]").forEach((el) => {
    if (el.style.transform) {
      el.style.transform = "none";
    }
  });

  const scaler = root.querySelector<HTMLElement>('[class*="scaler"]');
  const page = root.querySelector<HTMLElement>('[class*="page"]');

  if (scaler) {
    scaler.style.width = "297mm";
    scaler.style.height = "210mm";
    scaler.style.aspectRatio = "unset";
    scaler.style.position = "relative";
    scaler.style.overflow = "hidden";
  }

  if (page) {
    page.style.position = "absolute";
    page.style.top = "0";
    page.style.left = "0";
    page.style.width = "297mm";
    page.style.height = "210mm";
    page.style.transform = "none";
    page.style.boxShadow = "none";
  }
}
