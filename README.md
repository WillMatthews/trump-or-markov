# Trump or Markov Chain?

I have a guilty admission.
In this age of Large Language models, reasoning, generative AI, and all other marketing buzzwords, sometimes it's nice to take a step back.
I quite enjoy looking back at the simpler times, when the humble Markov Chain was the go-to tool for generating text.

Markov chains, for those who aren't in the know, are a very simple way of generating text - they're a type of stochastic model that describes a sequence of possible events in which the probability of each event depends only on the state attained in the previous event.
Cor - that's a mouthful, isn't it?

In layman's terms, a Markov chain is a way of predicting what word comes next in a sentence, based on the words that came before it.
You can train the model on a big chunk of text - all you do is see what words tend to come after other words, and then use that information to generate new text that sounds like it could have come from the original text.

For example, let's assume trained a Markov chain on the text "A B A B C A C".
If we start with the word "A", the Markov chain will predict that the next word is "B" 2/3 of the time and "C" 1/3 of the time.
If we start with the word "B", "A" comes next 1/2 of the time and "C" comes next 1/2 of the time. And so on.

It's very simple, and as a result generates text that is often nonsense.
Something else that's quite nonsensical is the political climate of the United States.
They also have a former president who is quite simple, and has a reasonable dataset of speeches and tweets.
It's a match made in heaven!

# Game Plan

The plan is simple.
We're going to train a Markov chain on a dataset of Donald Trump's speeches and tweets.
Then, we're going to generate some text using the Markov chain.
We're going to mix this generated text with some real Trump quotes, and we're going to see if you can tell the difference.
You probably won't be able to, because it's all nonsense anyway.

# WE'RE GONNA MAKE THE COMPUTER PAY FOR IT!

## Training the Markov Chain

I used n=2 for the Markov chain, which means that the model looks at the two previous words to predict the next word.
Information about the training corpus can be found in the `real-data` directory.
