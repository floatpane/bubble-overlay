import fs from "node:fs/promises";
import path from "node:path";
import matter from "gray-matter";
import { notFound } from "next/navigation";
import { MDXRemote } from "next-mdx-remote/rsc";
import remarkGfm from "remark-gfm";
import remarkFrontmatter from "remark-frontmatter";
import remarkCallout from "@r4ai/remark-callout";
import rehypeHighlight from "rehype-highlight";
import { Sidebar } from "@/components/Sidebar";

const CONTENT = path.join(process.cwd(), "content");

async function load(slug: string[]) {
	const file = path.join(CONTENT, `${slug.join("/")}.mdx`);
	try {
		return matter(await fs.readFile(file, "utf8"));
	} catch {
		return null;
	}
}

export default async function DocPage({
	params,
}: { params: Promise<{ slug: string[] }> }) {
	const { slug } = await params;
	const doc = await load(slug);
	if (!doc) notFound();

	return (
		<div className="docs-layout">
			<aside className="docs-sidebar">
				<Sidebar />
			</aside>
			<article className="docs-content">
				<MDXRemote
					source={doc.content}
					options={{
						mdxOptions: {
							remarkPlugins: [remarkFrontmatter, remarkGfm, remarkCallout],
							rehypePlugins: [rehypeHighlight],
						},
					}}
				/>
			</article>
		</div>
	);
}
