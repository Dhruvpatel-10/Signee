import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { ThemeProvider } from "@/components/Providers/ThemeProvides";
import { APP_NAME } from "@/config/constants";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  metadataBase: new URL("https://yourdomain.com"), 
  title: APP_NAME + " - SSH CA Cert Manager",
  description: "Easily manage and issue SSH CA certificates with a simple web app.",
  icons: "/logo/icon-signee.png",
  openGraph: {
    title: APP_NAME + " - SSH CA Cert Manager",
    description: "Manage SSH certificates easily with a secure, web-based tool.",
    url: "https://yourdomain.com",
    siteName: APP_NAME,
    images: [
      {
        url: "/logo/og-image.png", 
        width: 1200,
        height: 630,
        alt: APP_NAME + "Dashboard Preview",
      },
    ],
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: APP_NAME + "- SSH CA Cert Manager",
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
      <body className={`${geistSans.variable} ${geistMono.variable}`}>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          {children}
        </ThemeProvider>
      </body>
    </html>
  );
}
