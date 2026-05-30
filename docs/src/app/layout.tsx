import type { Metadata } from "next";
import "./globals.css";
import "highlight.js/styles/github-dark.css";

export const metadata: Metadata = {
	title: "bubble-overlay",
	description:
		"ANSI-aware overlay painter for Bubble Tea / lipgloss views. Paint modals, popups, and tooltips over a styled base.",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
	return (
		<html lang="en">
			<body>
				<header className="site-header">
					<a href="/" className="brand">
						bubble-overlay
					</a>
					<nav>
						<a href="https://github.com/floatpane/bubble-overlay">GitHub</a>
					</nav>
				</header>
				<main>{children}</main>
			</body>
		</html>
	);
}
