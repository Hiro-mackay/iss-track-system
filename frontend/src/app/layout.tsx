import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "ISS Tracker - Real-Time Position",
  description:
    "Real-time ISS tracking with orbital prediction and map visualization",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja" className="dark">
      <body>{children}</body>
    </html>
  );
}
