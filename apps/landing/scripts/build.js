import fs from 'node:fs';
import path from 'node:path';

const root = path.resolve(__dirname, '..');
const dist = path.join(root, 'dist');

const files = ['index.html', 'style.css', 'main.js'];

if (!fs.existsSync(dist)) {
  fs.mkdirSync(dist, { recursive: true });
}

for (const file of files) {
  const src = path.join(root, file);
  const dest = path.join(dist, file);
  if (fs.existsSync(src)) {
    fs.copyFileSync(src, dest);
    console.log(`Copied ${file} -> dist/${file}`);
  } else {
    console.warn(`Missing ${file}`);
  }
}

console.log('Landing page build complete.');
