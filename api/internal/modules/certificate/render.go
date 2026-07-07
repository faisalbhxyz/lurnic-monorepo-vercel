package certificate

import (
	"bytes"
	"dashlearn/internal/models"
	"fmt"
	"html/template"
	"strings"
	"time"
)

const TemplateMinarAcademy = "/templates/minar-academy"

type signerRenderBlock struct {
	SignatureImage string
	Name           string
	Role           string
	Org            string
}

type renderData struct {
	Title            string
	StudentName      string
	CourseName       string
	CertificateID    string
	IssuedDate       string
	PricingWord      string
	OrganizationName string
	CompletionLine   string
	BrandLogo        string
	WatermarkImage   string
	WatermarkOpacity float64
	DualSigners      bool
	PrimarySigner    signerRenderBlock
	SecondarySigner  signerRenderBlock
	ShowDownload     bool
}

func IsMinarTemplate(templatePath string) bool {
	path := strings.TrimSpace(templatePath)
	if path == "" || strings.HasPrefix(path, "/images/Certificat-") {
		return true
	}
	return path == TemplateMinarAcademy || strings.HasPrefix(path, TemplateMinarAcademy+"/")
}

func clampWatermarkOpacity(value uint8) uint8 {
	if value > 100 {
		return 100
	}
	return value
}

func watermarkOpacityValue(value uint8) float64 {
	return float64(clampWatermarkOpacity(value)) / 100
}

func renderCertificateHTML(cert models.StudentCertificate, showDownload bool) (string, error) {
	if !IsMinarTemplate(cert.TemplatePath) {
		return "", fmt.Errorf("unsupported certificate template: %s", cert.TemplatePath)
	}

	data := buildRenderData(cert, showDownload)
	var buf bytes.Buffer
	if err := minarTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func buildRenderData(cert models.StudentCertificate, showDownload bool) renderData {
	title := "Certificate of Completion"
	if cert.Title != nil && strings.TrimSpace(*cert.Title) != "" {
		title = strings.TrimSpace(*cert.Title)
	}

	completionLine := "has successfully completed"
	if cert.SubtitleOne != nil && strings.TrimSpace(*cert.SubtitleOne) != "" {
		completionLine = strings.TrimSpace(*cert.SubtitleOne)
	}

	orgName := "your organization"
	if cert.OrganizationName != nil && strings.TrimSpace(*cert.OrganizationName) != "" {
		orgName = strings.TrimSpace(*cert.OrganizationName)
	}

	pricingWord := "FREE"
	if cert.PricingModel == models.CoursePricingModelPaid {
		pricingWord = "PAID"
	}

	brandLogo := ""
	if cert.BrandLogo != nil {
		brandLogo = strings.TrimSpace(*cert.BrandLogo)
	}

	watermarkImage := ""
	if cert.WatermarkImage != nil {
		watermarkImage = strings.TrimSpace(*cert.WatermarkImage)
	}

	signatureImage := ""
	if cert.InstructorSignature != nil && strings.TrimSpace(*cert.InstructorSignature) != "" {
		signatureImage = strings.TrimSpace(*cert.InstructorSignature)
	} else if !cert.DualSignersEnabled && cert.OwnerSignature != nil && strings.TrimSpace(*cert.OwnerSignature) != "" {
		signatureImage = strings.TrimSpace(*cert.OwnerSignature)
	}

	signatureImage2 := ""
	if cert.OwnerSignature != nil && strings.TrimSpace(*cert.OwnerSignature) != "" {
		signatureImage2 = strings.TrimSpace(*cert.OwnerSignature)
	}

	signerName := ""
	if cert.SignerName != nil {
		signerName = strings.TrimSpace(*cert.SignerName)
	}

	signerRole := ""
	if cert.SignerRole != nil {
		signerRole = strings.TrimSpace(*cert.SignerRole)
	}

	signerOrg := ""
	if cert.SignerOrg != nil {
		signerOrg = strings.TrimSpace(*cert.SignerOrg)
	}

	signer2Name := ""
	if cert.Signer2Name != nil {
		signer2Name = strings.TrimSpace(*cert.Signer2Name)
	}

	signer2Role := ""
	if cert.Signer2Role != nil {
		signer2Role = strings.TrimSpace(*cert.Signer2Role)
	}

	signer2Org := ""
	if cert.Signer2Org != nil {
		signer2Org = strings.TrimSpace(*cert.Signer2Org)
	}

	return renderData{
		Title:            title,
		StudentName:      cert.StudentName,
		CourseName:       cert.CourseTitle,
		CertificateID:    cert.CertificateNumber,
		IssuedDate:       cert.IssuedAt.Format(issuedAtFormat),
		PricingWord:      pricingWord,
		OrganizationName: orgName,
		CompletionLine:   completionLine,
		BrandLogo:        brandLogo,
		WatermarkImage:   watermarkImage,
		WatermarkOpacity: watermarkOpacityValue(cert.WatermarkOpacity),
		DualSigners:      cert.DualSignersEnabled,
		PrimarySigner: signerRenderBlock{
			SignatureImage: signatureImage,
			Name:           signerName,
			Role:           signerRole,
			Org:            signerOrg,
		},
		SecondarySigner: signerRenderBlock{
			SignatureImage: signatureImage2,
			Name:           signer2Name,
			Role:           signer2Role,
			Org:            signer2Org,
		},
		ShowDownload: showDownload,
	}
}

var minarTemplate = template.Must(template.New("minar").Funcs(template.FuncMap{
	"safeURL": func(value string) template.URL {
		return template.URL(value)
	},
}).Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>{{ .Title }}</title>
<style>
  @page { size: A4 landscape; margin: 0; }
  html, body { margin: 0; padding: 0; height: 100%; }
  * { box-sizing: border-box; }
  body {
    font-family: 'Georgia', 'Times New Roman', serif;
    background: #e5e5e5;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
  }
/* Fixed A4 landscape canvas (297×210mm) — keep in sync with frontend/styles/minar-certificate.css */
  .scaler {
    width: 100%;
    aspect-ratio: 297 / 210;
    position: relative;
    overflow: hidden;
  }
  .page {
    position: absolute;
    top: 0;
    left: 0;
    width: 297mm;
    height: 210mm;
    overflow: hidden;
    background: #ffffff;
    box-shadow: 0 0 6px rgba(0,0,0,0.15);
    transform-origin: top left;
  }
  .download-btn {
    position: fixed; top: 16px; right: 16px; z-index: 20;
    border: none; border-radius: 8px; padding: 10px 14px;
    background: #0b463f; color: #fff; font-size: 14px; font-weight: 600;
    cursor: pointer; box-shadow: 0 6px 20px rgba(0,0,0,0.2);
  }
  .download-btn:hover { background: #0f5a51; }
  @media print {
    body { background: #ffffff; }
    .download-btn { display: none !important; }
    html, body { width: 297mm; height: 210mm; margin: 0; padding: 0; overflow: hidden; }
    .scaler { width: 297mm; height: 210mm; aspect-ratio: unset; overflow: hidden; }
    .page { width: 297mm; height: 210mm; transform: none !important; box-shadow: none; position: relative; top: 0; left: 0; }
    .footer-bar { bottom: 2.2mm; height: 4.2mm; }
  }
  .header-graphic { position: absolute; top: 0; left: 50%; transform: translateX(-50%); width: 100%; height: auto; display: block; z-index: 2; }
  .footer-bar { position: absolute; bottom: 1.8mm; left: 0; width: 100%; height: 4.8mm; z-index: 1; }
  .watermark-bg {
    position: absolute; top: 56%; left: 50%; transform: translate(-50%, -50%);
    width: 118mm; z-index: 1; object-fit: contain;
  }
  .content {
    position: absolute; top: 50mm; left: 50%; transform: translateX(-50%);
    z-index: 3; text-align: center; width: 100%;
    display: flex; flex-direction: column; align-items: center;
  }
  .brand-logo { display: block; width: 62mm; height: auto; margin-bottom: 2mm; }
  .brand-placeholder {
    width: 62mm; height: 22mm; margin-bottom: 2mm;
    border: 1px dashed #bbb; color: #888; font-size: 3.5mm;
    display: flex; align-items: center; justify-content: center;
  }
  .recipient-name {
    margin-top: 8mm; font-size: 34px; line-height: 1.15; color: #111; letter-spacing: 0.2px;
    font-weight: bold; text-transform: uppercase;
  }
  .completion-line { margin-top: 3mm; font-size: 6.3mm; line-height: 1.3; color: #555; }
  .course-name {
    margin-top: 5mm; font-size: 23px; line-height: 1.25; color: #111; font-weight: bold;
    max-width: 240mm; padding: 0 8mm;
  }
  .course-subline {
    margin-top: 2.5mm; font-size: 5.8mm; line-height: 1.35; color: #4a4a4a;
    max-width: 230mm; padding: 0 10mm;
  }
  .course-subline .pricing-word { color: #c53f47; font-weight: bold; }
  .meta-row {
    margin-top: 7mm; display: flex; justify-content: center; align-items: center;
    gap: 8mm; font-size: 4.9mm; color: #646464; width: auto;
  }
  .meta-row > div { width: 72mm; text-align: center; line-height: 1.35; word-break: break-word; }
  .meta-row > div:first-child { font-size: 16px; }
  .meta-row strong { color: #c53f47; font-weight: 600; }
  .meta-row .issued-at { color: #333; font-weight: 500; }
  .signatures {
    position: absolute; bottom: 14mm; left: 0; width: 100%;
    display: flex; justify-content: center; align-items: flex-end; z-index: 4;
  }
  .signatures.dual { justify-content: space-between; padding: 0 8%; }
  .sig { text-align: center; width: 72mm; }
  .signatures.dual .sig { width: 66mm; }
  .sig .name {
    font-size: 4.8mm; line-height: 1.25; color: #111; font-weight: 400; margin-top: 0.7mm;
  }
  .sig .sig-image {
    display: block; width: 100%; max-height: 16mm; object-fit: contain;
    object-position: bottom center; margin: 0 auto 1mm;
  }
  .sig .sig-line { width: 100%; height: 1px; background: #2c2c2c; margin-bottom: 2.2mm; }
  .sig .role { font-size: 4.8mm; line-height: 1.25; color: #333; margin-top: 0.7mm; }
  .sig .org { font-size: 4.5mm; line-height: 1.25; color: #555; margin-top: 0.5mm; }
</style>
</head>
<body>
{{ if .ShowDownload }}<button class="download-btn" onclick="window.print()">Download PDF</button>{{ end }}
<div class="scaler"><div class="page">
  <svg class="header-graphic" viewBox="0 0 1600 520" preserveAspectRatio="none" xmlns="http://www.w3.org/2000/svg">
    <defs>
      <linearGradient id="tealLeft" x1="0" y1="0" x2="800" y2="463" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#072e2a"/><stop offset="55%" stop-color="#0d5850"/><stop offset="100%" stop-color="#1c8577"/>
      </linearGradient>
      <linearGradient id="tealRight" x1="1600" y1="0" x2="800" y2="463" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#072e2a"/><stop offset="55%" stop-color="#0d5850"/><stop offset="100%" stop-color="#1c8577"/>
      </linearGradient>
      <linearGradient id="goldLeft" x1="0" y1="463" x2="800" y2="0" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#A87433"/><stop offset="45%" stop-color="#E6C88F"/><stop offset="60%" stop-color="#F3DEB4"/><stop offset="100%" stop-color="#C28F46"/>
      </linearGradient>
      <linearGradient id="goldRight" x1="1600" y1="463" x2="800" y2="0" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#A87433"/><stop offset="45%" stop-color="#E6C88F"/><stop offset="60%" stop-color="#F3DEB4"/><stop offset="100%" stop-color="#C28F46"/>
      </linearGradient>
      <radialGradient id="medalGrad" cx="50%" cy="50%" r="60%">
        <stop offset="0%" stop-color="#fff2c2"/><stop offset="60%" stop-color="#e0b04a"/><stop offset="100%" stop-color="#a97a1f"/>
      </radialGradient>
      <linearGradient id="tailGrad" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" stop-color="#f0c95c"/><stop offset="100%" stop-color="#b8860b"/>
      </linearGradient>
      <clipPath id="leftTriClip"><polygon points="0,0 800,0 0,463"/></clipPath>
      <clipPath id="rightTriClip"><polygon points="1600,0 800,0 1600,463"/></clipPath>
    </defs>
    <polygon points="0,0 800,0 0,463" fill="url(#tealLeft)"/>
    <polygon points="1600,0 800,0 1600,463" fill="url(#tealRight)"/>
    <polygon points="0,463 37.56,527.9 837.56,64.9 800,0" fill="url(#goldLeft)" stroke="#C28F46" stroke-width="60" stroke-linejoin="round" clip-path="url(#leftTriClip)"/>
    <polygon points="1600,463 1562.44,527.9 762.44,64.9 800,0" fill="url(#goldRight)" stroke="#C28F46" stroke-width="60" stroke-linejoin="round" clip-path="url(#rightTriClip)"/>
    <rect x="754" y="148" width="92" height="10" fill="#e0b04a"/>
    <polygon points="754,152 754,244 800,224 846,244 846,152" fill="url(#tailGrad)"/>
    <circle cx="800" cy="70" r="88" fill="url(#medalGrad)" stroke="#8a6316" stroke-width="4"/>
    <circle cx="800" cy="70" r="67" fill="none" stroke="#8a6316" stroke-width="2" opacity="0.55"/>
  </svg>
  <svg class="footer-bar" viewBox="0 0 2970 48" preserveAspectRatio="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    <defs>
      <linearGradient id="footerGrad" x1="0%" y1="0%" x2="100%" y2="0%">
        <stop offset="0%" stop-color="#0b463f"/><stop offset="50%" stop-color="#14685e"/><stop offset="100%" stop-color="#0b463f"/>
      </linearGradient>
    </defs>
    <rect x="0" y="0" width="2970" height="48" fill="url(#footerGrad)"/>
  </svg>
  {{ if .WatermarkImage }}
  <img class="watermark-bg" src="{{ .WatermarkImage | safeURL }}" alt="" style="opacity: {{ .WatermarkOpacity }};">
  {{ end }}
  <div class="content">
    {{ if .BrandLogo }}
    <img class="brand-logo" src="{{ .BrandLogo | safeURL }}" alt="Brand Logo">
    {{ else }}
    <div class="brand-placeholder">Your Logo</div>
    {{ end }}
    <div class="recipient-name">{{ .StudentName }}</div>
    <div class="completion-line">{{ .CompletionLine }}</div>
    <div class="course-name">{{ .CourseName }}</div>
    <div class="course-subline">a <span class="pricing-word">{{ .PricingWord }}</span> online course offered by {{ .OrganizationName }}.</div>
    <div class="meta-row">
      <div>CERTIFICATE ID: <strong>{{ .CertificateID }}</strong></div>
      <div>Issued <span class="issued-at">{{ .IssuedDate }}</span></div>
    </div>
  </div>
  <div class="signatures{{ if .DualSigners }} dual{{ end }}">
    <div class="sig">
      {{ if .PrimarySigner.SignatureImage }}<img class="sig-image" src="{{ .PrimarySigner.SignatureImage | safeURL }}" alt="Signature">{{ end }}
      <div class="sig-line"></div>
      {{ if .PrimarySigner.Name }}<div class="name">{{ .PrimarySigner.Name }}</div>{{ end }}
      {{ if .PrimarySigner.Role }}<div class="role">{{ .PrimarySigner.Role }}</div>{{ end }}
      {{ if .PrimarySigner.Org }}<div class="org">{{ .PrimarySigner.Org }}</div>{{ end }}
    </div>
    {{ if .DualSigners }}
    <div class="sig">
      {{ if .SecondarySigner.SignatureImage }}<img class="sig-image" src="{{ .SecondarySigner.SignatureImage | safeURL }}" alt="Signature">{{ end }}
      <div class="sig-line"></div>
      {{ if .SecondarySigner.Name }}<div class="name">{{ .SecondarySigner.Name }}</div>{{ end }}
      {{ if .SecondarySigner.Role }}<div class="role">{{ .SecondarySigner.Role }}</div>{{ end }}
      {{ if .SecondarySigner.Org }}<div class="org">{{ .SecondarySigner.Org }}</div>{{ end }}
    </div>
    {{ end }}
  </div>
</div></div>
</body>
</html>`))

func formatIssuedDatePreview(t time.Time) string {
	if t.IsZero() {
		return "[ISSUED_DATE]"
	}
	return t.Format(issuedAtFormat)
}
