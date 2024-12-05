class Coord
  attr_accessor :i, :j

  def initialize(i, j)
    @i, @j = i, j
  end

  def +(c)
    Coord.new(@i+c.i, @j+c.j)
  end

  def to_s
    "(#@i, #@j)"
  end

  def ==(other)
    @i == other.i && @j == other.j
  end

  alias eql? ==

  def hash
    [ @i, @j ].hash
  end
end

class DataPoint
  def initialize(matrix, point)
    @matrix, @point = matrix, point
  end

  def char
    @matrix[@point] || ?Z
  end

  def next(direction)
    DataPoint.new(@matrix, @point + direction)
  end

  def to_s
    "#@point=#{@matrix[@point]}"
  end
end

fname = ARGV[0] or raise 'Must supply a filename!'
lines = IO.readlines(fname, chomp: true)

@left   = Coord.new(0,  -1)
@right  = Coord.new(0,   1)
@up     = Coord.new(-1,  0)
@down   = Coord.new(1,   0)
@ludiag = @left + @up
@lddiag = @left + @down
@rudiag = @right + @up
@rddiag = @right + @down

@matrix = {}

(0...lines.size).each do |i|
  (0...lines[i].size).each do |j|
    @matrix[Coord.new(i, j)] = lines[i][j]
  end
end

@count = 0
@target = 'MAS'

def traverse(loc, direction, message)
  if message.empty?
    1
  elsif message[0] == loc.char
    traverse(loc.next(direction), direction, message[1..-1])
  else
    0
  end
end

def mas(loc, direction)
  traverse(loc, direction, @target)
end

def count_xmas(i, j)
  c = Coord.new(i, j)
  loc = DataPoint.new(@matrix, c)
  num = mas(loc, @rddiag) +
    mas(loc.next(@down+@down), @rudiag) +
    mas(loc.next(@right+@right), @lddiag) +
    mas(loc.next(@rddiag+@rddiag), @ludiag)

  if num > 1
    @count += 1
  end
end

(0...lines.size).each do |i|
  (0...lines[i].size).each do |j|
    count_xmas(i, j)
  end
end

puts @count
