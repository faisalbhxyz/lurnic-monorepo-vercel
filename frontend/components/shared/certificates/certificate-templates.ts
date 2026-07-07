export const CERTIFICATE_TEMPLATE_MINAR = "/templates/minar-academy";

export const DEFAULT_CERTIFICATE_TEMPLATE = CERTIFICATE_TEMPLATE_MINAR;

export function isMinarCertificateTemplate(templatePath?: string | null) {
  const path = templatePath?.trim() ?? "";
  if (!path || path.startsWith("/images/Certificat-")) {
    return true;
  }
  return (
    path === CERTIFICATE_TEMPLATE_MINAR ||
    path.startsWith(`${CERTIFICATE_TEMPLATE_MINAR}/`)
  );
}
