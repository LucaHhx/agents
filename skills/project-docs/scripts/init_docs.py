#!/usr/bin/env python3
"""
Initialize project documentation structure.

Usage:
    init_docs.py <project-root>

Creates docs/ directory with:
    project.md, specs.yaml, CHANGELOG.md, glossary.md, plans/
"""

import sys
from pathlib import Path
from datetime import date


def init_docs(project_root: str) -> bool:
    project_root = Path(project_root).resolve()
    docs_dir = project_root / "docs"

    if docs_dir.exists():
        print(f"Error: docs/ already exists at {docs_dir}")
        return False

    # Locate templates
    skill_root = Path(__file__).parent.parent
    templates_dir = skill_root / "assets" / "templates"

    if not templates_dir.exists():
        print(f"Error: Templates not found at {templates_dir}")
        return False

    # Create directory structure
    docs_dir.mkdir(parents=True)
    (docs_dir / "plans").mkdir()

    # Copy root templates with placeholder replacement
    today = date.today().isoformat()
    root_templates = ["project.md", "specs.yaml", "CHANGELOG.md", "glossary.md"]

    for template_name in root_templates:
        template_path = templates_dir / template_name
        if not template_path.exists():
            print(f"Warning: Template not found: {template_name}")
            continue
        content = template_path.read_text(encoding="utf-8")
        content = content.replace("{{DATE}}", today)
        (docs_dir / template_name).write_text(content, encoding="utf-8")

    print(f"Created docs/ at {docs_dir}")
    print("  - project.md")
    print("  - specs.yaml")
    print("  - CHANGELOG.md")
    print("  - glossary.md")
    print("  - plans/")
    return True


def main():
    if len(sys.argv) != 2:
        print("Usage: init_docs.py <project-root>")
        sys.exit(1)

    success = init_docs(sys.argv[1])
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
