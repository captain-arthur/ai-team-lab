#!/usr/bin/env python3
"""
Workflow initializer for the AI team workspace.

Given a task file in tasks/, creates a new project under projects/ with
phase directories and copies the relevant templates. Does not run AI;
it only scaffolds the workspace for each phase.

Usage:
    python scripts/run_workflow.py tasks/my-task.md
"""

import argparse
import shutil
import sys
from pathlib import Path


# Phase directories to create under the project
PHASES = [
    "01-manager",
    "02-research",
    "03-architecture",
    "04-engineering",
    "05-review",
    "06-documentation",
    "07-knowledge-extraction",
]

# Templates to copy: (source in templates/, destination path under phase dir)
TEMPLATE_COPIES = [
    ("research.md", "02-research", "research.md"),
    ("architecture.md", "03-architecture", "architecture.md"),
    ("final-report.md", "06-documentation", "final-report.md"),
    ("knowledge-extraction.md", "07-knowledge-extraction", "knowledge-extraction.md"),
]


def repo_root() -> Path:
    """Script lives in scripts/; repo root is parent of scripts/."""
    return Path(__file__).resolve().parent.parent


def project_name_from_task_path(task_path: Path) -> str:
    """Derive project name from task filename (e.g. my-task.md -> my-task)."""
    return task_path.stem


def run(task_file: str) -> None:
    task_path = Path(task_file).resolve()
    root = repo_root()
    tasks_dir = root / "tasks"
    projects_dir = root / "projects"
    templates_dir = root / "templates"

    # Validate task file
    if not task_path.is_file():
        print(f"Error: Task file not found: {task_path}", file=sys.stderr)
        sys.exit(1)
    try:
        task_path.relative_to(tasks_dir)
    except ValueError:
        if not task_path.exists():
            pass  # already handled above
        print(
            f"Warning: Task file is not under tasks/: {task_path}. Using filename for project name.",
            file=sys.stderr,
        )

    project_name = project_name_from_task_path(task_path)
    project_dir = projects_dir / project_name

    if project_dir.exists():
        print(f"Error: Project already exists: {project_dir}", file=sys.stderr)
        sys.exit(1)

    # Create phase directories
    project_dir.mkdir(parents=True)
    for phase in PHASES:
        (project_dir / phase).mkdir()

    # Copy templates into phase folders
    for src_name, phase_subdir, dest_name in TEMPLATE_COPIES:
        src = templates_dir / src_name
        dest = project_dir / phase_subdir / dest_name
        if src.is_file():
            shutil.copy2(src, dest)
        else:
            print(f"Warning: Template not found: {src}", file=sys.stderr)

    # Copy task file into project root for reference
    dest_intake = project_dir / "task-intake.md"
    shutil.copy2(task_path, dest_intake)

    # Minimal project README
    readme = project_dir / "README.md"
    readme.write_text(
        f"# Project: {project_name}\n\n"
        f"Initialized from task intake: `{task_path.name}`\n\n"
        f"See `task-intake.md` in this folder for the full intake. Run each phase using the prompts in `prompts/` and write outputs into the corresponding phase folder.\n"
    )

    print(f"Created project: {project_dir}")
    print(f"  Task intake copied to: {dest_intake}")
    print(f"  Phases: {', '.join(PHASES)}")
    print(f"\nNext: Run the Manager phase with `task-intake.md` as input, then work through each phase.")
    print(f"  Prompts: prompts/")
    print(f"  Workflow: WORKFLOW.md")


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Initialize a new AI team project from a task intake file."
    )
    parser.add_argument(
        "task_file",
        help="Path to task file (e.g. tasks/my-task.md)",
    )
    args = parser.parse_args()
    run(args.task_file)


if __name__ == "__main__":
    main()
