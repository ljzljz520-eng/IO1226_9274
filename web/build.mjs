import { mkdir, writeFile } from "node:fs/promises";

await mkdir("dist", { recursive: true });
const html = `<!doctype html><html lang="zh-CN"><head><meta charset="utf-8"><title>迷你弓手生存赛</title></head><body><main><h1>迷你弓手生存赛</h1><p>自动瞄准、升级与生存统计控制台</p></main></body></html>`;
await writeFile("dist/index.html", html);
