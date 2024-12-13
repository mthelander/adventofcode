def parseline(line)
  test, rest = line.split(?:)
  [ test, *rest.split(' ') ].map(&:to_i)
end

def operators(n)
  ops = %w[ * + ]
  args = [ ops ] * n
  ops.product(*args)
end

def testequation(eq, target)
  n1, op, n2, *rest = eq
  result = n1.method(op.to_sym).call(n2)
  if rest.empty?
    result == target
  else
    testequation([result, *rest], target)
  end
end

def istrue?(equation)
  test, *terms = equation
  ops = operators(terms.size-2)
  candidates = ops.map { |o| terms.zip(o).flatten.compact }
  candidates.any? { |eq| testequation(eq, test) }
end

filename = ARGV[0] or raise 'Must supply a filename!'
equations = IO.readlines(filename, chomp: true).map(&method(:parseline))
#p equations.select(&method(:istrue?)).map(&:first).reduce(&:+)
#equations.select(&method(:istrue?)).map(&:first).map(&method(:p))#.reduce(&:+)
equations.select(&method(:istrue?)).each { |x| p x }
