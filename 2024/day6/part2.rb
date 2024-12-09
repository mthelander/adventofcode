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

def explore(curr, direction, visited)
  v = Pair.new(curr, direction)
  return true if visited.key?(v)
  visited[v] = true

  ahead = curr+direction

  case @data[ahead]
  when nil
    false
  when ?#
    explore(curr, turn_right(direction), visited)
  else
    explore(ahead, direction, visited)
  end
end

def traverse(curr, direction, visited, path, obstructions)
  this = Pair.new(curr, direction)
  visited[this] = path[curr] = true

  right = turn_right(direction)
  ahead = curr+direction
  lookahead = @data[ahead]

  case lookahead
  when nil
    obstructions
  when ?#
    traverse(curr, right, visited, path, obstructions)
  else
    if lookahead == ?. && !path.key?(ahead)
      @data[ahead] = ?#
      if explore(curr, right, Hash[**visited])
        obstructions[ahead] = true
      end
      @data[ahead] = lookahead
    end
    traverse(ahead, direction, visited, path, obstructions)
  end
end

obstructions = traverse(@pos, UP, Hash[], Hash[], Hash[])
puts obstructions.keys.size
