#!/usr/bin/env python3
"""
Initialize a new full-stack project from templates.

Usage:
    init_project.py --name <project-name> --module <go-module-path> [--output <dir>]

Creates:
    {output}/{name}/server/   - Go backend
    {output}/{name}/web/      - Tauri + React frontend
"""

import argparse
import re
import shutil
import sys
from pathlib import Path

TEMPLATE_MODULE = "server"

TEXT_EXTENSIONS = {
    ".go", ".mod", ".yaml", ".yml", ".json", ".ts", ".tsx", ".js", ".jsx",
    ".css", ".html", ".toml", ".rs", ".md", ".gitignore", ".gitkeep",
}

TEXT_NAMES = {".env", ".env.production", ".gitignore", ".gitkeep"}


def is_text_file(path: Path) -> bool:
    return path.suffix in TEXT_EXTENSIONS or path.name in TEXT_NAMES


def title_case(kebab_name: str) -> str:
    return " ".join(word.capitalize() for word in kebab_name.split("-"))


def copy_and_replace(src_dir: Path, dst_dir: Path, replacements: dict, go_module: str = None):
    for src_file in sorted(src_dir.rglob("*")):
        if src_file.is_dir():
            continue
        rel = src_file.relative_to(src_dir)
        dst_file = dst_dir / rel
        dst_file.parent.mkdir(parents=True, exist_ok=True)

        if is_text_file(src_file):
            content = src_file.read_text(encoding="utf-8")
            # Go import path replacement
            if go_module:
                if src_file.suffix == ".go":
                    content = content.replace(f'"{TEMPLATE_MODULE}/', f'"{go_module}/')
                if src_file.name == "go.mod":
                    content = content.replace(
                        f"module {TEMPLATE_MODULE}", f"module {go_module}"
                    )
            # Generic placeholder replacement
            for placeholder, value in replacements.items():
                content = content.replace(placeholder, value)
            dst_file.write_text(content, encoding="utf-8")
        else:
            shutil.copy2(src_file, dst_file)


def main():
    parser = argparse.ArgumentParser(description="Initialize a new project from templates")
    parser.add_argument("--name", required=True, help="Project name (kebab-case)")
    parser.add_argument("--module", required=True, help="Go module path")
    parser.add_argument("--output", default=".", help="Parent directory (default: cwd)")
    args = parser.parse_args()

    # Validate name
    if not re.match(r"^[a-z][a-z0-9]*(-[a-z0-9]+)*$", args.name):
        print(f"Error: Project name must be kebab-case (e.g., 'my-app')")
        sys.exit(1)

    project_dir = Path(args.output).resolve() / args.name
    if project_dir.exists():
        print(f"Error: Directory already exists: {project_dir}")
        sys.exit(1)

    # Locate templates
    skill_root = Path(__file__).parent.parent
    templates_dir = skill_root / "assets" / "templates"
    if not templates_dir.exists():
        print(f"Error: Templates not found at {templates_dir}")
        sys.exit(1)

    # Build replacements
    project_title = title_case(args.name)
    app_identifier = f"com.hz.{args.name}"
    replacements = {
        "{{PROJECT_NAME}}": args.name,
        "{{PROJECT_TITLE}}": project_title,
        "{{GO_MODULE}}": args.module,
        "{{APP_IDENTIFIER}}": app_identifier,
    }

    # Copy server template
    print(f"Creating {project_dir}/server/ ...")
    copy_and_replace(
        templates_dir / "server",
        project_dir / "server",
        replacements,
        go_module=args.module,
    )

    # Copy web template
    print(f"Creating {project_dir}/web/ ...")
    copy_and_replace(
        templates_dir / "web",
        project_dir / "web",
        replacements,
    )

    print(f"\nProject '{args.name}' created at {project_dir}")
    print(f"\nNext steps:")
    print(f"  1. Initialize docs:  use project-docs skill or run init_docs.py")
    print(f"  2. Backend setup:    cd {project_dir}/server && go mod tidy")
    print(f"  3. Frontend setup:   cd {project_dir}/web && npm install")
    print(f"  4. Start backend:    cd {project_dir}/server && go run main.go")
    print(f"  5. Start frontend:   cd {project_dir}/web && npm run tauri:dev")


if __name__ == "__main__":
    main()
