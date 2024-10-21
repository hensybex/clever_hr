import os

def process_file(file_path, base_dir, file_extension, subdirectory):
    with open(file_path, 'r') as file:
        lines = file.readlines()

    # Handle Go files
    if file_extension == ".go":
        # Remove everything before the "package" string
        package_index = next((i for i, line in enumerate(lines) if 'package' in line), None)
        if package_index is not None:
            lines = lines[package_index:]

    # Prepare the new comment with the subdirectory
    relative_path = os.path.relpath(file_path, base_dir)
    new_comment = f"// {subdirectory}/{relative_path}\n\n"

    # Check if the comment already exists and update it if necessary
    if lines and lines[0].startswith("// "):
        existing_comment = lines[0].strip()
        if existing_comment != new_comment.strip():
            lines[0] = new_comment
    else:
        lines.insert(0, new_comment)

    # Ensure exactly two newlines after the comment
    if len(lines) > 2 and lines[1].strip() == "" and lines[2].strip() == "":
        lines = lines[:2] + [line for line in lines[2:] if line.strip() != ""]
    elif len(lines) > 1 and lines[1].strip() == "":
        lines.insert(2, "\n")

    with open(file_path, 'w') as file:
        file.writelines(lines)

def process_directory(root_dir, extension, subdirectory):
    for subdir, _, files in os.walk(root_dir):
        for file in files:
            if file.endswith(extension) and not file.endswith(".g.dart"):
                file_path = os.path.join(subdir, file)
                process_file(file_path, root_dir, extension, subdirectory)

def run_for_all_targets(targets):
    for target in targets:
        print(f"Processing {target['name']} files...")
        process_directory(target['root_directory'], target['file_extension'], target['subdirectory'])

if __name__ == "__main__":
    # List of directories and their respective file extensions to process
    targets = [
        {"name": "API", "root_directory": "api/internal", "file_extension": ".go", "subdirectory": "internal"},
        {"name": "APP", "root_directory": "app/lib", "file_extension": ".dart", "subdirectory": "lib"}
    ]

    run_for_all_targets(targets)
