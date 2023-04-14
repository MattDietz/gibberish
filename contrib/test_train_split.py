import random
import sys


def main():
    seed = 42
    random.seed(42)

    path = sys.argv[1]
    with open(path, "r") as f:
        lines = f.readlines()
    random.shuffle(lines)
    train = lines[:int(len(lines) * 0.8)]
    test = lines[int(len(lines) * 0.8):]
    with open(sys.argv[2], "w") as f:
        f.writelines(train)
    with open(sys.argv[3], "w") as f:
        f.writelines(test)


if __name__ == "__main__":
    main()
