import os

def process_file(file_path, base_dir):
    with open(file_path, 'r') as file:
        lines = file.readlines()

    # Remove everything before the "package" string
    package_index = next((i for i, line in enumerate(lines) if 'package' in line), None)
    if package_index is not None:
        lines = lines[package_index:]

    # Insert the new comment at the beginning
    relative_path = os.path.relpath(file_path, base_dir)
    comment = f"// {relative_path}\n\n"
    lines.insert(0, comment)

    with open(file_path, 'w') as file:
        file.writelines(lines)

def process_directory(root_dir, extension):
    for subdir, _, files in os.walk(root_dir):
        for file in files:
            if file.endswith(extension):
                file_path = os.path.join(subdir, file)
                process_file(file_path, root_dir)

def run_for_all_targets(targets):
    for target in targets:
        print(f"Processing {target['name']} files...")
        process_directory(target['root_directory'], target['file_extension'])

if __name__ == "__main__":
    # List of directories and their respective file extensions to process
    targets = [
        {"name": "API", "root_directory": "api/internal", "file_extension": ".go"},
        {"name": "APP", "root_directory": "app/lib", "file_extension": ".dart"}
    ]

    run_for_all_targets(targets)
