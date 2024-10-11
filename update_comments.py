import os

def process_file(file_path):
    with open(file_path, 'r') as file:
        lines = file.readlines()

    # Remove everything before the "package" string
    package_index = next((i for i, line in enumerate(lines) if 'package' in line), None)
    if package_index is not None:
        lines = lines[package_index:]

    # Insert the new comment at the beginning
    relative_path = os.path.relpath(file_path, 'api')
    comment = f"// {relative_path}\n\n"
    lines.insert(0, comment)

    with open(file_path, 'w') as file:
        file.writelines(lines)

def process_directory(root_dir):
    for subdir, _, files in os.walk(root_dir):
        for file in files:
            if file.endswith('.go'):
                file_path = os.path.join(subdir, file)
                process_file(file_path)

if __name__ == "__main__":
    root_directory = 'api/internal'
    process_directory(root_directory)
