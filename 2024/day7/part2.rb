class Integer
  def concat(other)
    (self * (10 ** other.numdigits)) + other
  end

  def numdigits
    Math.log10(self).to_i + 1
  end
end

OPERATORS = %i[ * + concat ]
@pm = Hash[]

def parseline(line)
  test, rest = line.split(?:)
  [ test, *rest.split(' ') ].map(&:to_i)
end

def solve?(target, n, *equation)
  return n == target if equation.empty?
  return false if n > target

  n2, *rest = equation

  OPERATORS.each do |op|
    result = n.method(op).call(n2)
    if solve?(target, result, *rest)
      return true
    end
  end
  return false
end

filename = ARGV[0] or raise 'Must supply a filename!'
equations = IO.readlines(filename, chomp: true).map(&method(:parseline))
p equations.filter { |eq| solve?(*eq) }.map(&:first).reduce(&:+)
