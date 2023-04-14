import random
import sys

import gibberish
import numpy as np

vowels = 'aeiou'
consonants = 'bcdfghjklmnpqrstvwxyz'
everything = ''.join([chr(c) for c in range(33, 64)] + [chr(c) for c in range(97, 126)])

mean_word_length = 4.5
mean_sentence_length = 17.5
std_word_length = 2



def garbage(base_str):
    return ''.join(random.choice(base_str)
                   for _ in range(int(np.random.normal(mean_word_length,
                                                       std_word_length))))

def sentence_len():
    return int(np.random.exponential(mean_sentence_length))

def sentence_range():
    return range(int(np.random.exponential(sentence_len())))


if __name__ == '__main__':
    if len(sys.argv) != 2:
        print('Usage: python gibberish.py <num_tokens>')
        sys.exit(1)

    to_generate = int(sys.argv[1])
    g = gibberish.Gibberish()
    for w in g.generate_words(to_generate):
        print(w)

    for i in range(to_generate):
        print(' '.join(w for w in g.generate_words(sentence_len())))

    for i in range(to_generate):
        print(garbage(vowels))

    for i in range(to_generate):
        print(' '.join(garbage(vowels) for _ in sentence_range()))

    for i in range(to_generate):
        print(garbage(consonants))

    for i in range(to_generate):
        print(' '.join(garbage(consonants) for _ in sentence_range()))

    for i in range(to_generate):
        print(garbage(everything))

    for i in range(to_generate):
        print(' '.join(garbage(everything) for _ in sentence_range()))
