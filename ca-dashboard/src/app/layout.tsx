import type { Metadata } from "next";
import "./globals.css";
import { ThemeProvider } from "@/components/Providers/ThemeProvides";
import { APP_NAME, MY_DOMAIN } from "@/config/constants";
import { Toaster } from "@/components/ui/sonner";

export const metadata: Metadata = {
  metadataBase: new URL(MY_DOMAIN), 
  title: APP_NAME + "- SSH CA Cert Manager",
  description: "Easily manage and issue SSH CA certificates with a simple web app.",
  icons: "/logo/icon-signee.png",
  openGraph: {
    title: APP_NAME + " - SSH CA Cert Manager",
    description: "Manage SSH certificates easily with a secure, web-based tool.",
    url: MY_DOMAIN,
    siteName: "Signee",
    images: [
      {
        url: "/logo/og-image.png", 
        width: 1200,
        height: 630,
        alt: APP_NAME + " Dashboard Preview",
      },
    ],
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: APP_NAME + " - SSH CA Cert Manager",
    description: "Manage SSH certificates easily with a secure, web-based tool.",
    images: ["/logo/og-image.png"],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          {children}
          <Toaster position="top-right" />
        </ThemeProvider>
      </body>
    </html>
  );
}
