Pair = Struct.new(:x, :y) do
  def +(other)
    Pair.new(x+other.x, y+other.y)
  end
end

LEFT   = Pair.new(0,  -1)
RIGHT  = Pair.new(0,   1)
UP     = Pair.new(-1,  0)
DOWN   = Pair.new(1,   0)

DIRECTIONS = [ UP, RIGHT, DOWN, LEFT ]

def build_data(filename)
  {}.tap do |data|
    lines = IO.readlines(filename, chomp: true)
    (0...lines.size).each do |i|
      (0...lines[i].size).each do |j|
        data[Pair.new(i, j)] = lines[i][j]
      end
    end
  end
end

fname = ARGV[0] or raise 'Must supply a filename!'
@data = build_data(fname)
@pos  = @data.key(?^)

def turn_right(dir)
  DIRECTIONS[(DIRECTIONS.index(dir) + 1) % 4]
end

def traverse(curr, direction, visited)
  visited.add(curr)

  right = turn_right(direction)
  ahead = curr+direction
  lookahead = @data[ahead]

  case lookahead
  when nil
    visited
  when ?#
    traverse(curr, right, visited)
  else
    traverse(ahead, direction, visited)
  end
end

path = traverse(@pos, UP, Set[])
puts path.size
