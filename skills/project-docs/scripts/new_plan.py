#!/usr/bin/env python3
"""
Create a new plan directory.

Usage:
    new_plan.py <plans-directory> <plan-name>

Creates <plans-directory>/<plan-name>/ with:
    plan.md, tasks.md, decisions.md, changelog.md, testing.md, testing/
"""

import sys
import re
from pathlib import Path
from datetime import date


def new_plan(plans_dir: str, plan_name: str) -> bool:
    plans_dir = Path(plans_dir).resolve()

    # Validate plan name (kebab-case)
    if not re.match(r"^[a-z0-9]+(-[a-z0-9]+)*$", plan_name):
        print(f"Error: Plan name must be kebab-case (e.g., 'user-auth', 'api-redesign')")
        return False

    plan_dir = plans_dir / plan_name

    if plan_dir.exists():
        print(f"Error: Plan already exists at {plan_dir}")
        return False

    if not plans_dir.exists():
        print(f"Error: Plans directory not found: {plans_dir}")
        print("Run init_docs.py first to create the docs structure.")
        return False

    # Locate templates
    skill_root = Path(__file__).parent.parent
    templates_dir = skill_root / "assets" / "templates"

    # Create plan directory and testing/ subdirectory
    plan_dir.mkdir(parents=True)
    (plan_dir / "testing").mkdir()

    # Copy plan templates
    today = date.today().isoformat()
    plan_templates = {
        "plan.md": "plan.md",
        "tasks.md": "tasks.md",
        "decisions.md": "decisions.md",
        "changelog.md": "plan_changelog.md",
        "testing.md": "testing.md",
    }

    for target_name, template_name in plan_templates.items():
        template_path = templates_dir / template_name
        if not template_path.exists():
            print(f"Warning: Template not found: {template_name}")
            continue
        content = template_path.read_text(encoding="utf-8")
        content = content.replace("{{DATE}}", today)
        content = content.replace("{{PLAN_NAME}}", plan_name)
        (plan_dir / target_name).write_text(content, encoding="utf-8")

    print(f"Created plan at {plan_dir}")
    print("  - plan.md")
    print("  - tasks.md")
    print("  - decisions.md")
    print("  - changelog.md")
    print("  - testing.md")
    print("  - testing/")
    return True


def main():
    if len(sys.argv) != 3:
        print("Usage: new_plan.py <plans-directory> <plan-name>")
        sys.exit(1)

    success = new_plan(sys.argv[1], sys.argv[2])
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
